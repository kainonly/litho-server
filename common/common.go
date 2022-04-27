package common

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/nats-io/nats.go"
	"github.com/tencentyun/cos-go-sdk-v5"
	"github.com/weplanx/go/encryption"
	"github.com/weplanx/go/engine"
	"github.com/weplanx/go/passport"
	openapi "github.com/weplanx/openapi/client"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"time"
)

const TokenClaimsKey = "token-claims"

var (
	AuthExpired  = errors.New("认证已失效，令牌超出有效期")
	AuthConflict = errors.New("认证已失效，已被新终端占用")
)

type Inject struct {
	Values      *Values
	MongoClient *mongo.Client
	Db          *mongo.Database
	Redis       *redis.Client
	Nats        *nats.Conn
	Js          nats.JetStreamContext
	Open        *openapi.OpenAPI
	Cipher      *encryption.Cipher
	HID         *encryption.HID
	Cos         *cos.Client
}

func SetValues(path string) (values *Values, err error) {
	if _, err = os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("静态配置不存在，请检查路径 [%s]", path)
	}
	var config []byte
	if config, err = ioutil.ReadFile(path); err != nil {
		return
	}
	if err = yaml.Unmarshal(config, &values); err != nil {
		return
	}
	return
}

type Values struct {
	TrustedProxies []string                 `yaml:"trusted_proxies"`
	Name           string                   `yaml:"name"`
	Namespace      string                   `yaml:"namespace"`
	Key            string                   `yaml:"key"`
	Cors           Cors                     `yaml:"cors"`
	Database       Database                 `yaml:"database"`
	Redis          Redis                    `yaml:"redis"`
	Nats           Nats                     `yaml:"nats"`
	Engines        map[string]engine.Option `yaml:"engines"`
	OpenAPI        OpenAPI                  `yaml:"openapi"`
	Passport       passport.Option          `yaml:"passport"`
	QCloud         QCloud                   `yaml:"qcloud"`
	Feishu         Feishu                   `yaml:"feishu"`
	Email          Email                    `yaml:"email"`
}

func (x *Values) KeyName(v ...string) string {
	return fmt.Sprintf(`%s:%s`, x.Namespace, strings.Join(v, ":"))
}

func (x *Values) EventName(v string) string {
	return fmt.Sprintf(`%s.events.%s`, x.Namespace, v)
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

type Nats struct {
	Hosts []string `yaml:"hosts"`
	Nkey  string   `yaml:"nkey"`
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

type OpenAPI struct {
	Url    string `yaml:"url"`
	Key    string `yaml:"key"`
	Secret string `yaml:"secret"`
}

type Feishu struct {
	Redirect          string `yaml:"redirect"`
	AppId             string `yaml:"app_id"`
	AppSecret         string `yaml:"app_secret"`
	EncryptKey        string `yaml:"encrypt_key"`
	VerificationToken string `yaml:"verification_token"`
}

type Email struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type Subscriptions struct {
	*sync.Map
}

func BoolToP(v bool) *bool {
	return &v
}

func ObjectIDToP(v interface{}) *primitive.ObjectID {
	if id, ok := v.(primitive.ObjectID); ok {
		return &id
	}
	return nil
}

type LoginLogDto struct {
	Time      time.Time          `bson:"time"`
	V         string             `bson:"v"`
	User      primitive.ObjectID `bson:"user"`
	Username  string             `bson:"username"`
	Email     string             `bson:"email"`
	TokenId   string             `bson:"token_id"`
	Ip        string             `bson:"ip"`
	Detail    bson.M             `bson:"detail"`
	UserAgent string             `bson:"user_agent"`
}

func NewLoginLogV10(data User, jti string, ip string, agent string) *LoginLogDto {
	return &LoginLogDto{
		Time:      time.Now(),
		V:         "v1.0",
		User:      data.ID,
		Username:  data.Username,
		Email:     data.Email,
		TokenId:   jti,
		Ip:        ip,
		UserAgent: agent,
	}
}
