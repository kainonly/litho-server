package bootstrap

import (
	"context"
	"database/sql"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/config"
	"github.com/go-playground/validator/v10"
	"github.com/hertz-contrib/binding/go_playground"
	"github.com/hertz-contrib/requestid"
	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"
	"github.com/weplanx/go/captcha"
	"github.com/weplanx/go/cipher"
	"github.com/weplanx/go/csrf"
	"github.com/weplanx/go/help"
	"github.com/weplanx/go/locker"
	"github.com/weplanx/go/passport"
	"gopkg.in/yaml.v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
	"regexp"
	"server/common"
	"strings"
	"time"
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
			DSN:                  v.Database.Url,
			PreferSimpleProtocol: true,
		}),
		&gorm.Config{
			Logger:                 log,
			SkipDefaultTransaction: true,
			PrepareStmt:            true,
		},
	); err != nil {
		return
	}
	var db *sql.DB
	if db, err = orm.DB(); err != nil {
		return
	}
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)
	db.SetConnMaxLifetime(time.Hour)
	return
}

func UseRedis(v *common.Values) (client *redis.Client, err error) {
	opts, err := redis.ParseURL(v.Database.Redis)
	if err != nil {
		return
	}
	client = redis.NewClient(opts)
	if err = client.Ping(context.TODO()).Err(); err != nil {
		return
	}
	return
}

func UseNats(v *common.Values) (nc *nats.Conn, err error) {
	if nc, err = nats.Connect(
		strings.Join(v.Nats.Hosts, ","),
		nats.MaxReconnects(-1),
		nats.Token(v.Nats.Token),
	); err != nil {
		return
	}
	return
}

func UseJetStream(nc *nats.Conn) (nats.JetStreamContext, error) {
	return nc.JetStream(nats.PublishAsyncMaxPending(256))
}

func UseKeyValue(v *common.Values, js nats.JetStreamContext) (nats.KeyValue, error) {
	return js.CreateKeyValue(&nats.KeyValueConfig{Bucket: v.Namespace})
}

func UsePassport(v *common.Values) *passport.Passport {
	return passport.New(
		passport.SetKey(v.Key),
		passport.SetIssuer(v.Domain),
	)
}

func UseCsrf(v *common.Values) *csrf.Csrf {
	return csrf.New(
		csrf.SetKey(v.Key),
		csrf.SetDomain(v.Domain),
	)
}

func UseCipher(v *common.Values) (*cipher.Cipher, error) {
	return cipher.New(v.Key)
}

func UseLocker(client *redis.Client) *locker.Locker {
	return locker.New(client)
}

func UseCaptcha(client *redis.Client) *captcha.Captcha {
	return captcha.New(client)
}

func UseHertz(v *common.Values) (h *server.Hertz, err error) {
	if v.Address == "" {
		return
	}
	vd := go_playground.NewValidator()
	vd.SetValidateTag("vd")
	vdx := vd.Engine().(*validator.Validate)
	vdx.RegisterValidation("snake", func(fl validator.FieldLevel) bool {
		matched, errX := regexp.MatchString("^[a-z_]+$", fl.Field().Interface().(string))
		if errX != nil {
			return false
		}
		return matched
	})
	vdx.RegisterValidation("sort", func(fl validator.FieldLevel) bool {
		matched, errX := regexp.MatchString("^[a-z_.]+:(-1|1)$", fl.Field().Interface().(string))
		if errX != nil {
			return false
		}
		return matched
	})

	opts := []config.Option{
		server.WithHostPorts(v.Address),
		server.WithCustomValidator(vd),
	}

	if os.Getenv("MODE") != "release" {
		opts = append(opts, server.WithExitWaitTime(0))
	}

	opts = append(opts)
	h = server.Default(opts...)
	h.Use(
		help.ErrorHandler(),
		requestid.New(),
	)

	return
}
