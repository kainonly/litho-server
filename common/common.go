package common

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/tencentyun/cos-go-sdk-v5"
	"github.com/weplanx/go/encryption"
	"github.com/weplanx/go/engine"
	"github.com/weplanx/go/passport"
	"go.mongodb.org/mongo-driver/mongo"
	"strings"
)

const TokenClaimsKey = "token-claims"

type Inject struct {
	Values      *Values
	MongoClient *mongo.Client
	Db          *mongo.Database
	Redis       *redis.Client
	Passport    *passport.Passport
	Cipher      *encryption.Cipher
	Idx         *encryption.IDx
	Cos         *cos.Client
}

type Values struct {
	Address        string                   `yaml:"address"`
	TrustedProxies []string                 `yaml:"trusted_proxies"`
	Name           string                   `yaml:"name"`
	Key            string                   `yaml:"key"`
	Cors           Cors                     `yaml:"cors"`
	Database       Database                 `yaml:"database"`
	Redis          Redis                    `yaml:"redis"`
	Pulsar         Pulsar                   `yaml:"pulsar"`
	Passport       passport.Option          `yaml:"passport"`
	Engines        map[string]engine.Option `yaml:"engines"`
	QCloud         QCloud                   `yaml:"qcloud"`
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

type Pulsar struct {
	Url    string            `yaml:"url"`
	Token  string            `yaml:"token"`
	Topics map[string]string `yaml:"topics"`
}

type QCloud struct {
	SecretID  string    `yaml:"secret_id"`
	SecretKey string    `yaml:"secret_key"`
	Cos       QCloudCos `yaml:"cos"`
}

type QCloudCos struct {
	Bucket  string `yaml:"bucket"`
	Region  string `yaml:"region"`
	Expired int64  `yaml:"expired"`
}

func (x *Values) KeyName(v ...string) string {
	return fmt.Sprintf(`%s:%s`, x.Name, strings.Join(v, ":"))
}
