package bootstrap

import (
	"context"
	"fmt"
	"github.com/caarlos0/env/v9"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/config"
	"github.com/hertz-contrib/requestid"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nkeys"
	"github.com/redis/go-redis/v9"
	"github.com/weplanx/go-wpx/captcha"
	"github.com/weplanx/go-wpx/cipher"
	"github.com/weplanx/go-wpx/csrf"
	"github.com/weplanx/go-wpx/locker"
	"github.com/weplanx/go-wpx/passport"
	"github.com/weplanx/go-wpx/sessions"
	"github.com/weplanx/go-wpx/values"
	"github.com/weplanx/server/common"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"os"
	"strings"
	"time"
)

func LoadStaticValues() (v *common.Values, err error) {
	v = new(common.Values)
	if err = env.Parse(v); err != nil {
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
	return &values.Service{
		KeyValue: kv,
		Cipher:   cipher,
		Values:   v.DynamicValues,
	}
}

func UseCsrf(v *common.Values) *csrf.Csrf {
	return csrf.New(
		csrf.SetKey(v.Key),
	)
}

func UseCipher(v *common.Values) (*cipher.Cipher, error) {
	return cipher.New(v.Key)
}

func UseSessions(v *common.Values, rdb *redis.Client) *sessions.Service {
	return &sessions.Service{
		Namespace: v.Namespace,
		Redis:     rdb,
		Values:    v.DynamicValues,
	}
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

func UseHertz(v *common.Values) (h *server.Hertz, err error) {
	opts := []config.Option{
		server.WithHostPorts(v.Address),
	}

	if os.Getenv("MODE") != "release" {
		opts = append(opts, server.WithExitWaitTime(0))
	}

	opts = append(opts)

	h = server.Default(opts...)

	h.Use(
		requestid.New(),
	)

	return
}
