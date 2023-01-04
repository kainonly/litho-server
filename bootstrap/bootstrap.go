package bootstrap

import (
	"context"
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/config"
	"github.com/go-redis/redis/v8"
	"github.com/go-resty/resty/v2"
	"github.com/hertz-contrib/requestid"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nkeys"
	"github.com/weplanx/server/api"
	"github.com/weplanx/server/common"
	"github.com/weplanx/transfer"
	"github.com/weplanx/utils/captcha"
	"github.com/weplanx/utils/csrf"
	"github.com/weplanx/utils/dsl"
	"github.com/weplanx/utils/kv"
	"github.com/weplanx/utils/locker"
	"github.com/weplanx/utils/passport"
	"github.com/weplanx/utils/sessions"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"os"
	"strings"
	"time"
)

func LoadStaticValues() (values *common.Values, err error) {
	values = new(common.Values)
	if err = env.Parse(values); err != nil {
		return
	}
	values.DynamicValues = &kv.DEFAULT
	values.DynamicValues.DSL = map[string]*kv.DSLOption{
		"users": {
			Keys: []string{"_id", "email", "name", "avatar", "status", "sessions", "last", "create_time", "update_time"},
		},
	}
	return
}

// UseMongoDB
// https://www.mongodb.com/docs/drivers/go/current/
// https://pkg.go.dev/go.mongodb.org/mongo-driver/mongo
func UseMongoDB(values *common.Values) (*mongo.Client, error) {
	return mongo.Connect(
		context.TODO(),
		options.Client().ApplyURI(values.Database.Mongo),
	)
}

// UseDatabase
// https://www.mongodb.com/docs/drivers/go/current/
// https://pkg.go.dev/go.mongodb.org/mongo-driver/mongo
func UseDatabase(values *common.Values, client *mongo.Client) (db *mongo.Database) {
	option := options.Database().
		SetWriteConcern(writeconcern.New(writeconcern.WMajority()))
	return client.Database(values.Database.Name, option)
}

// UseRedis
// https://github.com/go-redis/redis
func UseRedis(values *common.Values) (client *redis.Client, err error) {
	opts, err := redis.ParseURL(values.Database.Redis)
	if err != nil {
		return
	}
	client = redis.NewClient(opts)
	if err = client.Ping(context.TODO()).Err(); err != nil {
		return
	}
	return
}

// UseNats
// https://docs.nats.io/using-nats/developer
// https://github.com/nats-io/nats.go
func UseNats(values *common.Values) (nc *nats.Conn, err error) {
	var kp nkeys.KeyPair
	if kp, err = nkeys.FromSeed([]byte(values.Nats.Nkey)); err != nil {
		return
	}
	defer kp.Wipe()
	var pub string
	if pub, err = kp.PublicKey(); err != nil {
		return
	}
	if !nkeys.IsValidPublicUserKey(pub) {
		return nil, fmt.Errorf("nkey 验证失败")
	}
	if nc, err = nats.Connect(
		strings.Join(values.Nats.Hosts, ","),
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

// UseJetStream
// https://docs.nats.io/using-nats/developer/develop_jetstream
func UseJetStream(nc *nats.Conn) (nats.JetStreamContext, error) {
	return nc.JetStream(nats.PublishAsyncMaxPending(256))
}

// UseKeyValue
// https://docs.nats.io/using-nats/developer/develop_jetstream/kv
func UseKeyValue(values *common.Values, js nats.JetStreamContext) (nats.KeyValue, error) {
	return js.CreateKeyValue(&nats.KeyValueConfig{Bucket: values.Namespace})
}

func UseKV(values *common.Values, keyvalue nats.KeyValue) *kv.KV {
	return kv.New(
		kv.SetNamespace(values.Namespace),
		kv.SetKeyValue(keyvalue),
		kv.SetDynamicValues(values.DynamicValues),
	)
}

func UseSessions(values *common.Values, redis *redis.Client) *sessions.Sessions {
	return sessions.New(
		sessions.SetNamespace(values.Namespace),
		sessions.SetRedis(redis),
		sessions.SetDynamicValues(values.DynamicValues),
	)
}

func UseDSL(values *common.Values, db *mongo.Database) (*dsl.DSL, error) {
	return dsl.New(
		dsl.SetNamespace(values.Namespace),
		dsl.SetDatabase(db),
		dsl.SetDynamicValues(values.DynamicValues),
	)
}

func UseCsrf(values *common.Values) *csrf.Csrf {
	return csrf.New(
		csrf.SetKey(values.Key),
	)
}

func UsePassport(values *common.Values) *passport.Passport {
	return passport.New(
		passport.SetNamespace(values.Namespace),
		passport.SetKey(values.Key),
	)
}

func UseLocker(values *common.Values, client *redis.Client) *locker.Locker {
	return locker.New(
		locker.SetNamespace(values.Namespace),
		locker.SetRedis(client),
	)
}

func UseCaptcha(values *common.Values, client *redis.Client) *captcha.Captcha {
	return captcha.New(
		captcha.SetNamespace(values.Namespace),
		captcha.SetRedis(client),
	)
}

// UseTransfer
// https://github.com/weplanx/transfer
func UseTransfer(values *common.Values, db *mongo.Database, js nats.JetStreamContext) (*transfer.Transfer, error) {
	return transfer.New(
		transfer.SetNamespace(values.Namespace),
		transfer.SetDatabase(db),
		transfer.SetJetStream(js),
	)
}

// UseHttpClients 创建请求客户端
func UseHttpClients() *common.HttpClients {
	return &common.HttpClients{
		Feishu: resty.New().
			SetBaseURL("https://open.feishu.cn/open-apis"),
	}
}

// UseHertz
// https://www.cloudwego.io/zh/docs/hertz/reference/config
func UseHertz(values *common.Values) (h *server.Hertz, err error) {
	opts := []config.Option{
		server.WithHostPorts(values.Address),
	}

	if os.Getenv("MODE") != "release" {
		opts = append(opts, server.WithExitWaitTime(0))
	}

	h = server.Default(opts...)

	h.Use(
		requestid.New(),
	)

	return
}

// UseTest
func UseTest() (api *api.API, err error) {
	values, err := LoadStaticValues()
	if err != nil {
		panic(err)
	}
	if api, err = NewAPI(values); err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if _, err = api.Initialize(ctx); err != nil {
		return
	}

	return
}
