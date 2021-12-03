package common

import (
	"github.com/go-redis/redis/v8"
	"github.com/weplanx/go/passport"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/fx"
)

const TokenClaimsKey = "token-claims"

type Inject struct {
	fx.In

	Values      *Values
	MongoClient *mongo.Client
	Db          *mongo.Database
	Redis       *redis.Client
	Passport    *passport.Passport
}

type Values struct {
	Address  string          `yaml:"address"`
	Name     string          `yaml:"name"`
	Key      string          `yaml:"key"`
	Cors     Cors            `yaml:"cors"`
	Database Database        `yaml:"database"`
	Redis    Redis           `yaml:"redis"`
	Passport passport.Option `yaml:"passport"`
}

type Cors struct {
	AllowOrigins     []string `yaml:"allowOrigins"`
	AllowMethods     []string `yaml:"allowMethods"`
	AllowHeaders     []string `yaml:"allowHeaders"`
	ExposeHeaders    []string `yaml:"exposeHeaders"`
	AllowCredentials bool     `yaml:"allowCredentials"`
	MaxAge           int      `yaml:"maxAge"`
}

type Database struct {
	Uri    string `yaml:"uri"`
	DbName string `yaml:"dbName"`
}

type Redis struct {
	Uri string `yaml:"uri"`
}

func (x *Values) RedisKey(name string) string {
	return x.Name + ":" + name
}
