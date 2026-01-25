package bootstrap

import (
	"context"
	"database/sql"
	"os"
	"server/common"
	"strings"

	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/config"
	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/hertz-contrib/cors"
	"github.com/kainonly/go/captcha"
	"github.com/kainonly/go/cipher"
	"github.com/kainonly/go/csrf"
	"github.com/kainonly/go/help"
	"github.com/kainonly/go/locker"
	"github.com/kainonly/go/passport"
	"github.com/kainonly/go/vd"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/redis/go-redis/v9"
	"gopkg.in/yaml.v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func LoadStaticValues(path string) (v *common.Values, err error) {
	v = new(common.Values)
	var b []byte
	if b, err = os.ReadFile(path); err != nil {
		return
	}
	if err = yaml.Unmarshal(b, &v); err != nil {
		return
	}
	return
}

func UseGorm(v *common.Values) (orm *gorm.DB, err error) {
	var log logger.Interface
	if v.IsSqlDebug() {
		log = logger.Default.LogMode(logger.Info)
	}
	if orm, err = gorm.Open(
		postgres.New(postgres.Config{
			DSN:                  v.Database.DSN,
			PreferSimpleProtocol: true,
		}),
		&gorm.Config{
			Logger:                 log,
			SkipDefaultTransaction: v.Database.Gorm.SkipDefaultTransaction,
			PrepareStmt:            v.Database.Gorm.PrepareStmt,
		},
	); err != nil {
		return
	}
	var db *sql.DB
	if db, err = orm.DB(); err != nil {
		return
	}
	db.SetMaxIdleConns(v.Database.Pool.MaxIdleConns)
	db.SetMaxOpenConns(v.Database.Pool.MaxOpenConns)
	db.SetConnMaxLifetime(v.Database.Pool.ConnMaxLife)
	return
}

func UseRedis(v *common.Values) (client *redis.Client, err error) {
	opts, err := redis.ParseURL(v.Database.Redis)
	if err != nil {
		return
	}
	opts.Protocol = 2
	client = redis.NewClient(opts)
	if err = client.Ping(context.TODO()).Err(); err != nil {
		return
	}
	return
}

func UseNats(v *common.Values) (nc *nats.Conn, err error) {
	if nc, err = nats.Connect(
		strings.Join(v.Nats.Hosts, ","),
		nats.Token(v.Nats.Token),
		nats.RetryOnFailedConnect(true),
		nats.MaxReconnects(-1),
	); err != nil {
		return
	}
	return
}

func UseJetStream(nc *nats.Conn) (jetstream.JetStream, error) {
	return jetstream.New(nc)
}

func UsePassport(v *common.Values) *passport.Passport {
	return passport.New(
		passport.SetKey(v.App.Key),
		passport.SetIssuer(v.App.Namespace),
	)
}

func UseCsrf(v *common.Values) *csrf.Csrf {
	return csrf.New(
		csrf.SetKey(v.App.Key),
		csrf.SetDomain(v.Network.Domain),
	)
}

func UseCipher(v *common.Values) (*cipher.Cipher, error) {
	return cipher.New(v.App.Key)
}

func UseLocker(client *redis.Client) *locker.Locker {
	return locker.New(client)
}

func UseCaptcha(client *redis.Client) *captcha.Captcha {
	return captcha.New(client)
}

func UseHertz(v *common.Values) (h *server.Hertz, err error) {
	if v.App.Address == "" {
		return
	}

	vdx := vd.New()
	opts := []config.Option{
		server.WithHostPorts(v.App.Address),
		server.WithCustomValidatorFunc(func(request *protocol.Request, i interface{}) error {
			return vdx.Validate(request)
		}),
	}

	opts = append(opts)
	h = server.Default(opts...)
	h.Use(
		help.ErrorHandler(),
		cors.New(cors.Config{
			AllowOrigins:     v.Cors.Origins,
			AllowMethods:     v.Cors.Methods,
			AllowHeaders:     v.Cors.Headers,
			ExposeHeaders:    v.Cors.ExposeHeaders,
			AllowCredentials: v.Cors.AllowCredentials,
			MaxAge:           v.Cors.MaxAge,
		}),
	)
	return
}
