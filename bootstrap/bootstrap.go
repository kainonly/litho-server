package bootstrap

import (
	"context"
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/config"
	"github.com/go-redis/redis/v8"
	"github.com/hertz-contrib/cors"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nkeys"
	"github.com/weplanx/server/api"
	"github.com/weplanx/server/common"
	"github.com/weplanx/transfer"
	"github.com/weplanx/utils/captcha"
	"github.com/weplanx/utils/dsl"
	"github.com/weplanx/utils/kv"
	"github.com/weplanx/utils/locker"
	"github.com/weplanx/utils/passport"
	"github.com/weplanx/utils/sessions"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"strings"
	"time"
)

// LoadStaticValues 加载静态配置
func LoadStaticValues() (values *common.Values, err error) {
	values = new(common.Values)
	if err = env.Parse(values); err != nil {
		return
	}
	values.DynamicValues = &kv.DEFAULT
	return
}

// UseDatabase 初始化数据库
// 配置文档 https://gorm.io/zh_CN/
func UseDatabase(values *common.Values) (db *gorm.DB, err error) {
	return gorm.Open(postgres.New(postgres.Config{
		DSN:                  values.Gorm,
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		SkipDefaultTransaction:                   true,
		DisableForeignKeyConstraintWhenMigrating: true,
	})
}

// UseRedis 初始化 Redis
// 配置文档 https://github.com/go-redis/redis
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

// UseNats 初始化 Nats
// 配置文档 https://docs.nats.io/using-nats/developer
// SDK https://github.com/nats-io/nats.go
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

// UseJetStream 初始化流
// 说明 https://docs.nats.io/using-nats/developer/develop_jetstream
func UseJetStream(nc *nats.Conn) (nats.JetStreamContext, error) {
	return nc.JetStream(nats.PublishAsyncMaxPending(256))
}

// UseKeyValue 初始 KeyValue
// 说明 https://docs.nats.io/using-nats/developer/develop_jetstream/kv
func UseKeyValue(values *common.Values, js nats.JetStreamContext) (nats.KeyValue, error) {
	return js.CreateKeyValue(&nats.KeyValueConfig{Bucket: values.Namespace})
}

// UseKV 使用动态配置
func UseKV(values *common.Values, keyvalue nats.KeyValue) *kv.KV {
	return kv.New(
		kv.SetNamespace(values.Namespace),
		kv.SetKeyValue(keyvalue),
		kv.SetDynamicValues(values.DynamicValues),
	)
}

// UseSessions 使用会话
func UseSessions(values *common.Values, redis *redis.Client) *sessions.Sessions {
	return sessions.New(
		sessions.SetNamespace(values.Namespace),
		sessions.SetRedis(redis),
		sessions.SetDynamicValues(values.DynamicValues),
	)
}

// UseDSL 使用通用查询
func UseDSL(values *common.Values, db *mongo.Database) (*dsl.DSL, error) {
	return dsl.New(
		dsl.SetNamespace(values.Namespace),
		dsl.SetDatabase(db),
		dsl.SetDynamicValues(values.DynamicValues),
	)
}

// UsePassport 使用鉴权
func UsePassport(values *common.Values) *passport.Passport {
	return passport.New(
		passport.SetNamespace(values.Namespace),
		passport.SetKey(values.Key),
	)
}

// UseLocker 使用锁定
func UseLocker(values *common.Values, client *redis.Client) *locker.Locker {
	return locker.New(
		locker.SetNamespace(values.Namespace),
		locker.SetRedis(client),
	)
}

// UseCaptcha 使用验证
func UseCaptcha(values *common.Values, client *redis.Client) *captcha.Captcha {
	return captcha.New(
		captcha.SetNamespace(values.Namespace),
		captcha.SetRedis(client),
	)
}

// UseTransfer 初始日志传输
// https://github.com/weplanx/transfer
func UseTransfer(values *common.Values, db *mongo.Database, js nats.JetStreamContext) (*transfer.Transfer, error) {
	return transfer.New(
		transfer.SetNamespace(values.Namespace),
		transfer.SetDatabase(db),
		transfer.SetJetStream(js),
	)
}

// UseHertz 使用 Hertz
// 配置文档 https://www.cloudwego.io/zh/docs/hertz/reference/config
func UseHertz(values *common.Values) (h *server.Hertz, err error) {
	opts := []config.Option{
		server.WithHostPorts(values.Address),
	}

	if os.Getenv("MODE") != "release" {
		opts = append(opts, server.WithExitWaitTime(0))
	}

	h = server.Default(opts...)

	// 全局中间件
	h.Use(cors.New(cors.Config{
		AllowOrigins:     values.Hosts,
		AllowMethods:     []string{"GET", "POST", "PATCH", "PUT", "DELETE"},
		AllowHeaders:     []string{"Content-Type", "X-Pagesize", "X-Page"},
		AllowCredentials: true,
		ExposeHeaders:    []string{"X-Total"},
		MaxAge:           time.Hour * 2,
	}))

	return
}

// UseTest 初始测试
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
