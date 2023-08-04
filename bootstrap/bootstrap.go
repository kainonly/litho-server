package bootstrap

import (
	"context"
	"fmt"
	"github.com/caarlos0/env/v9"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/config"
	"github.com/hertz-contrib/obs-opentelemetry/provider"
	"github.com/hertz-contrib/obs-opentelemetry/tracing"
	"github.com/hertz-contrib/requestid"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nkeys"
	"github.com/redis/go-redis/v9"
	"github.com/weplanx/go/captcha"
	"github.com/weplanx/go/cipher"
	"github.com/weplanx/go/csrf"
	"github.com/weplanx/go/locker"
	"github.com/weplanx/go/passport"
	"github.com/weplanx/go/rest"
	"github.com/weplanx/go/sessions"
	"github.com/weplanx/go/values"
	"github.com/weplanx/server/common"
	"github.com/weplanx/transfer"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"gopkg.in/yaml.v3"
	"os"
	"reflect"
	"strings"
	"time"
)

func LoadStaticValues() (v *common.Values, err error) {
	v = new(common.Values)
	if err = env.Parse(v); err != nil {
		return
	}
	var b []byte
	if b, err = os.ReadFile("./config/default.values.yml"); err != nil {
		return
	}
	if err = yaml.Unmarshal(b, &v.Extra); err != nil {
		return
	}
	return
}

func UseMongoDB(v *common.Values) (*mongo.Client, error) {
	return mongo.Connect(
		context.TODO(),
		options.Client().ApplyURI(v.Database.Url),
	)
}

func UseDatabase(v *common.Values, client *mongo.Client) (db *mongo.Database) {
	option := options.Database().
		SetWriteConcern(writeconcern.Majority())
	return client.Database(v.Database.Name, option)
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
	var kp nkeys.KeyPair
	if kp, err = nkeys.FromSeed([]byte(v.Nats.Nkey)); err != nil {
		return
	}
	defer kp.Wipe()
	var pub string
	if pub, err = kp.PublicKey(); err != nil {
		return
	}
	if !nkeys.IsValidPublicUserKey(pub) {
		return nil, fmt.Errorf("nkey fail")
	}
	if nc, err = nats.Connect(
		strings.Join(v.Nats.Hosts, ","),
		nats.MaxReconnects(5),
		nats.ReconnectWait(2*time.Second),
		nats.ReconnectJitter(500*time.Millisecond, 2*time.Second),
		nats.Nkey(pub, func(nonce []byte) ([]byte, error) {
			sig, _ := kp.Sign(nonce)
			return sig, nil
		}),
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

func UseValues(v *common.Values, kv nats.KeyValue, cipher *cipher.Cipher) *values.Service {
	return values.New(
		values.SetNamespace(v.Namespace),
		values.SetKeyValue(kv),
		values.SetCipher(cipher),
		values.SetType(reflect.TypeOf(common.Extra{})),
	)
}

func UseSessions(v *common.Values, rdb *redis.Client) *sessions.Service {
	return sessions.New(
		sessions.SetNamespace(v.Namespace),
		sessions.SetRedis(rdb),
		sessions.SetDynamicValues(v.DynamicValues),
	)
}

func UseRest(
	v *common.Values,
	mgo *mongo.Client,
	db *mongo.Database,
	rdb *redis.Client,
	js nats.JetStreamContext,
	keyvalue nats.KeyValue,
) *rest.Service {
	return rest.New(
		rest.SetNamespace(v.Namespace),
		rest.SetMongoClient(mgo),
		rest.SetDatabase(db),
		rest.SetRedis(rdb),
		rest.SetJetStream(js),
		rest.SetKeyValue(keyvalue),
		rest.SetDynamicValues(v.DynamicValues),
	)
}

func UseCsrf(v *common.Values) *csrf.Csrf {
	return csrf.New(
		csrf.SetKey(v.Key),
	)
}

func UseCipher(v *common.Values) (*cipher.Cipher, error) {
	return cipher.New(v.Key)
}

func UsePassport(v *common.Values) *passport.Passport {
	return passport.New(
		passport.SetNamespace(v.Namespace),
		passport.SetKey(v.Key),
	)
}

func UseLocker(v *common.Values, client *redis.Client) *locker.Locker {
	return locker.New(
		locker.SetNamespace(v.Namespace),
		locker.SetRedis(client),
	)
}

func UseCaptcha(v *common.Values, client *redis.Client) *captcha.Captcha {
	return captcha.New(
		captcha.SetNamespace(v.Namespace),
		captcha.SetRedis(client),
	)
}

func UseTransfer(v *common.Values, js nats.JetStreamContext) (*transfer.Transfer, error) {
	return transfer.New(
		transfer.SetNamespace(v.Namespace),
		transfer.SetJetStream(js),
	)
}

func ProviderOpenTelemetry(v *common.Values) provider.OtelProvider {
	return provider.NewOpenTelemetryProvider(
		provider.WithServiceName(v.Namespace),
		provider.WithExportEndpoint(v.Otlp.Endpoint),
		provider.WithInsecure(),
	)
}

func UseHertz(v *common.Values) (h *server.Hertz, err error) {
	tracer, cfg := tracing.NewServerTracer()
	opts := []config.Option{
		server.WithHostPorts(v.Address),
		tracer,
	}
	if os.Getenv("MODE") != "release" {
		opts = append(opts, server.WithExitWaitTime(0))
	}
	opts = append(opts)
	h = server.Default(opts...)
	h.Use(
		requestid.New(),
		tracing.ServerMiddleware(cfg),
	)
	return
}
