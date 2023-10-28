package bootstrap

import (
	"context"
	"fmt"
	"github.com/caarlos0/env/v10"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/config"
	"github.com/hertz-contrib/obs-opentelemetry/provider"
	"github.com/hertz-contrib/obs-opentelemetry/tracing"
	"github.com/hertz-contrib/requestid"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nkeys"
	"github.com/redis/go-redis/v9"
	transfer "github.com/weplanx/collector/client"
	"github.com/weplanx/go/captcha"
	"github.com/weplanx/go/cipher"
	"github.com/weplanx/go/csrf"
	"github.com/weplanx/go/help"
	"github.com/weplanx/go/locker"
	"github.com/weplanx/go/passport"
	"github.com/weplanx/go/rest"
	"github.com/weplanx/go/sessions"
	"github.com/weplanx/go/values"
	"github.com/weplanx/server/common"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"gopkg.in/yaml.v3"
	"os"
	"reflect"
	"strings"
)

func LoadStaticValues(path string) (v *common.Values, err error) {
	v = new(common.Values)
	if err = env.Parse(v); err != nil {
		return
	}
	var b []byte
	if b, err = os.ReadFile(path); err != nil {
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
		nats.MaxReconnects(-1),
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

func UseValues(kv nats.KeyValue, cipher *cipher.Cipher) *values.Service {
	return values.New(
		values.SetKeyValue(kv),
		values.SetCipher(cipher),
		values.SetType(reflect.TypeOf(common.Extra{})),
	)
}

func UseSessions(v *common.Values, rdb *redis.Client) *sessions.Service {
	return sessions.New(
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
	xcipher *cipher.Cipher,
) *rest.Service {
	return rest.New(
		rest.SetMongoClient(mgo),
		rest.SetDatabase(db),
		rest.SetRedis(rdb),
		rest.SetJetStream(js),
		rest.SetKeyValue(keyvalue),
		rest.SetDynamicValues(v.DynamicValues),
		rest.SetCipher(xcipher),
	)
}

func UseCsrf(v *common.Values) *csrf.Csrf {
	return csrf.New(
		csrf.SetKey(v.Key),
		csrf.SetDomain(v.XDomain),
	)
}

func UseCipher(v *common.Values) (*cipher.Cipher, error) {
	return cipher.New(v.Key)
}

func UseAPIPassport(v *common.Values) *common.APIPassport {
	return passport.New(
		passport.SetIssuer(v.Namespace),
		passport.SetKey(v.Key),
	)
}

func UseLocker(client *redis.Client) *locker.Locker {
	return locker.New(client)
}

func UseCaptcha(client *redis.Client) *captcha.Captcha {
	return captcha.New(client)
}

func UseTransfer(js nats.JetStreamContext) (*transfer.Client, error) {
	return transfer.New(js)
}

func UseInflux(v *common.Values) influxdb2.Client {
	return influxdb2.NewClient(v.Influx.Url, v.Influx.Token)
}

func ProviderOpenTelemetry(v *common.Values) provider.OtelProvider {
	return provider.NewOpenTelemetryProvider(
		provider.WithServiceName(v.Namespace),
		provider.WithExportEndpoint(v.Otlp.Endpoint),
		provider.WithDeploymentEnvironment(v.Mode),
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

	tracer, cfg := tracing.NewServerTracer()
	opts := []config.Option{
		server.WithHostPorts(v.Address),
		server.WithCustomValidator(help.Validator()),
		tracer,
	}

	if os.Getenv("MODE") != "release" {
		opts = append(opts, server.WithExitWaitTime(0))
	}
	opts = append(opts)
	h = server.Default(opts...)
	h.Use(
		requestid.New(),
		help.EHandler(),
		tracing.ServerMiddleware(cfg),
	)

	return
}
