package bootstrap

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/config"
	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
	"github.com/hertz-contrib/cors"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nkeys"
	"github.com/weplanx/server/common"
	"github.com/weplanx/server/common/captcha"
	"github.com/weplanx/server/common/locker"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

var Provides = wire.NewSet(
	LoadStaticValues,
	UseMongoDB,
	UseDatabase,
	UseRedis,
	UseNats,
	UseJetStream,
	UseHertz,
	wire.Struct(new(captcha.Captcha), "*"),
	wire.Struct(new(locker.Locker), "*"),
)

// LoadStaticValues 加载静态配置
// 默认配置路径 ./config/config.yml
func LoadStaticValues() (values *common.Values, err error) {
	path := "./config/config.yml"
	if _, err = os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("静态配置不存在，请检查路径 [%s]", path)
	}
	var b []byte
	if b, err = ioutil.ReadFile(path); err != nil {
		return
	}
	if err = yaml.Unmarshal(b, &values); err != nil {
		return
	}
	return
}

// UseMongoDB 初始化 MongoDB
// 配置文档 https://www.mongodb.com/docs/drivers/go/current/
// https://pkg.go.dev/go.mongodb.org/mongo-driver/mongo
func UseMongoDB(values *common.Values) (*mongo.Client, error) {
	return mongo.Connect(
		context.TODO(),
		options.Client().ApplyURI(values.Database.Uri),
	)
}

// UseDatabase 初始化数据库
// 配置文档 https://www.mongodb.com/docs/drivers/go/current/
// https://pkg.go.dev/go.mongodb.org/mongo-driver/mongo
func UseDatabase(client *mongo.Client, values *common.Values) (db *mongo.Database) {
	option := options.Database().
		SetWriteConcern(writeconcern.New(writeconcern.WMajority()))
	return client.Database(values.Database.Db, option)
}

// UseRedis 初始化 Redis
// 配置文档 https://github.com/go-redis/redis
func UseRedis(values *common.Values) (client *redis.Client, err error) {
	opts, err := redis.ParseURL(values.Redis.Uri)
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

// UseHertz 使用 Hertz
// 配置文档 https://www.cloudwego.io/zh/docs/hertz/reference/config
func UseHertz(values *common.Values) (h *server.Hertz, err error) {
	opts := []config.Option{
		server.WithHostPorts(":3000"),
	}

	if os.Getenv("MODE") != "release" {
		opts = append(opts, server.WithExitWaitTime(0))
	}

	h = server.Default(opts...)

	// 全局中间件
	h.Use(cors.New(cors.Config{
		AllowOrigins:     values.Cors.AllowOrigins,
		AllowMethods:     values.Cors.AllowMethods,
		AllowHeaders:     values.Cors.AllowHeaders,
		AllowCredentials: values.Cors.AllowCredentials,
		ExposeHeaders:    values.Cors.ExposeHeaders,
		MaxAge:           time.Duration(values.Cors.MaxAge) * time.Second,
	}))

	return
}
