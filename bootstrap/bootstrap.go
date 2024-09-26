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
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nkeys"
	"github.com/redis/go-redis/v9"
	"github.com/weplanx/go/captcha"
	"github.com/weplanx/go/cipher"
	"github.com/weplanx/go/csrf"
	"github.com/weplanx/go/help"
	"github.com/weplanx/go/locker"
	"github.com/weplanx/go/passport"
	"github.com/weplanx/server/common"
	"os"
	"strings"
)

func LoadStaticValues(path string) (v *common.Values, err error) {
	v = new(common.Values)
	if err = env.Parse(v); err != nil {
		return
	}
	//var b []byte
	//if b, err = os.ReadFile(path); err != nil {
	//	return
	//}
	//if err = yaml.Unmarshal(b, &v.Extra); err != nil {
	//	return
	//}
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

func UseCsrf(v *common.Values) *csrf.Csrf {
	return csrf.New(
		csrf.SetKey(v.Key),
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

	opts := []config.Option{
		server.WithHostPorts(v.Address),
		server.WithCustomValidator(help.Validator()),
	}

	var tracer config.Option
	var tracerCfg *tracing.Config
	if *v.Otlp.Enabled {
		tracer, tracerCfg = tracing.NewServerTracer()
		opts = append(opts, tracer)
	}

	if os.Getenv("MODE") != "release" {
		opts = append(opts, server.WithExitWaitTime(0))
	}

	opts = append(opts)
	h = server.Default(opts...)
	h.Use(
		requestid.New(),
		help.EHandler(),
	)

	if tracerCfg != nil {
		h.Use(tracing.ServerMiddleware(tracerCfg))
	}

	return
}
