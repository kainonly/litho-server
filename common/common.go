package common

import (
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/go-redis/redis/v8"
	"github.com/nats-io/nats.go"
	"github.com/weplanx/utils/kv"
	"github.com/weplanx/utils/passport"
	"go.mongodb.org/mongo-driver/mongo"
	"strings"
)

type Inject struct {
	Values    *Values
	Mongo     *mongo.Client
	Db        *mongo.Database
	Redis     *redis.Client
	Nats      *nats.Conn
	JetStream nats.JetStreamContext
	KeyValue  nats.KeyValue
}

type Values struct {
	// 监听地址
	Address string `env:"ADDRESS" envDefault:":3000"`

	// 命名空间
	Namespace string `env:"NAMESPACE,required"`

	// 密钥
	Key string `env:"KEY,required"`

	// 跨域
	Hosts []string `env:"HOSTS" envSeparator:","`

	// 数据库
	Database `envPrefix:"DATABASE_"`

	// NATS 配置
	Nats `envPrefix:"NATS_"`

	// 动态配置
	*kv.DynamicValues `env:"-"`
}

type Database struct {
	// MongoDB 连接 Uri
	Mongo string `env:"MONGO,required"`

	// MongoDB 数据库名称
	Name string `env:"NAME,required"`

	// Redis 连接 Uri
	Redis string `env:"REDIS,required"`
}

type Nats struct {
	// Nats 连接地址
	Hosts []string `env:"HOSTS,required" envSeparator:","`

	// Nats Nkey 认证
	Nkey string `env:"NKEY,required"`
}

// Name 生成空间名称
func (x Values) Name(v ...string) string {
	return fmt.Sprintf(`%s:%s`, x.Namespace, strings.Join(v, ":"))
}

// GetClaims 获取授权标识
func GetClaims(c *app.RequestContext) (claims passport.Claims) {
	value, ok := c.Get("identity")
	if !ok {
		return
	}
	return value.(passport.Claims)
}
