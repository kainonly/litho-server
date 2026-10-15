package bootstrap

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"litho-api/common"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/config"
	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/hertz-contrib/cors"
	"github.com/hertz-contrib/obs-opentelemetry/provider"
	"github.com/hertz-contrib/obs-opentelemetry/tracing"
	"github.com/kainonly/collector/v3/transfer"
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
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/writeconcern"
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
	if err = yaml.Unmarshal(b, v); err != nil {
		return
	}
	if b, err = os.ReadFile(filepath.Join(filepath.Dir(path), "resources.yml")); err != nil {
		return
	}
	decoder := yaml.NewDecoder(bytes.NewReader(b))
	decoder.KnownFields(true)
	if err = decoder.Decode(&v.Resources); err != nil {
		return
	}
	if err = v.Resources.Prepare(); err != nil {
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
	var opts *redis.Options
	if opts, err = redis.ParseURL(v.Database.Redis); err != nil {
		return
	}
	opts.Protocol = 2
	client = redis.NewClient(opts)
	if err = client.Ping(context.TODO()).Err(); err != nil {
		return
	}
	return
}

func UseMongo(v *common.Values) (*mongo.Client, error) {
	return mongo.Connect(
		options.Client().ApplyURI(v.Database.Mongodb.Url),
	)
}

func UseMongoDatabase(v *common.Values, client *mongo.Client) (db *mongo.Database) {
	option := options.Database().
		SetWriteConcern(writeconcern.Majority())
	return client.Database(v.Database.Mongodb.Database, option)
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

func UseTransfer(v *common.Values, nc *nats.Conn) (*transfer.Transfer, error) {
	return transfer.New(context.TODO(), v.Namespace, nc)
}

func ProviderOpenTelemetry(v *common.Values) provider.OtelProvider {
	return provider.NewOpenTelemetryProvider(
		provider.WithServiceName(v.Otlp.Name),
		provider.WithExportEndpoint(v.Otlp.Endpoint),
		provider.WithDeploymentEnvironment(v.Otlp.Environment),
		provider.WithHeaders(map[string]string{
			"Authorization": fmt.Sprintf(`Bearer %s`, v.Otlp.Token),
		}),
		provider.WithEnableTracing(true),
		provider.WithEnableMetrics(true),
		provider.WithEnableCompression(),
		provider.WithInsecure(),
	)
}

func UseHertz(v *common.Values) (h *server.Hertz, err error) {
	if v.Address == "" {
		return
	}

	vdx := vd.Default()
	opts := []config.Option{
		server.WithHostPorts(v.Address),
		server.WithCustomValidatorFunc(func(_ *protocol.Request, req any) error {
			return vdx.Validate(req)
		}),
	}

	var tracer config.Option
	var tracerCfg *tracing.Config
	tracer, tracerCfg = tracing.NewServerTracer()
	opts = append(opts, tracer)

	h = server.Default(opts...)
	h.Use(
		help.ErrorHandler(),
		cors.New(cors.Config{
			AllowOrigins: v.Cors.Origins,
			AllowMethods: []string{"GET", "POST"},
			AllowHeaders: []string{"Origin", "Content-Length", "Content-Type",
				"X-Page", "X-Pagesize", "X-Route", "X-Requested-With",
			},
			ExposeHeaders:    []string{"X-Total"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		}),
	)
	if tracerCfg != nil {
		h.Use(tracing.ServerMiddleware(tracerCfg))
	}

	return
}
