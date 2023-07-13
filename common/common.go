package common

import (
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"
	"github.com/weplanx/go/captcha"
	"github.com/weplanx/go/cipher"
	"github.com/weplanx/go/locker"
	"github.com/weplanx/go/passport"
	"github.com/weplanx/go/values"
	"go.mongodb.org/mongo-driver/mongo"
	"strings"
	"time"
)

type Inject struct {
	V         *Values
	Mgo       *mongo.Client
	Db        *mongo.Database
	RDb       *redis.Client
	JetStream nats.JetStreamContext
	KeyValue  nats.KeyValue
	Cipher    *cipher.Cipher
	Passport  *passport.Passport
	Captcha   *captcha.Captcha
	Locker    *locker.Locker
}

type Values struct {
	Address   string `env:"ADDRESS" envDefault:":3000"`
	Namespace string `env:"NAMESPACE,required"`
	Key       string `env:"KEY,required"`
	Cors      `envPrefix:"CORS_"`
	Database  `envPrefix:"DATABASE_"`
	Nats      `envPrefix:"NATS_"`

	*values.DynamicValues
}

type Cors struct {
	AllowOrigins     []string      `env:"ALLOW_ORIGINS,required" envSeparator:","`
	AllowMethods     []string      `env:"ALLOW_METHODS,required" envSeparator:","`
	AllowHeaders     []string      `env:"ALLOW_HEADERS,required" envSeparator:","`
	ExposeHeaders    []string      `env:"EXPOSE_HEADERS,required" envSeparator:","`
	AllowCredentials bool          `env:"ALLOW_CREDENTIALS,required"`
	MaxAge           time.Duration `env:"MAX_AGE,required"`
}

type Database struct {
	Url   string `env:"URL,required"`
	Name  string `env:"NAME,required"`
	Redis string `env:"REDIS,required"`
}

type Nats struct {
	Hosts []string `env:"HOSTS,required" envSeparator:","`
	Nkey  string   `env:"NKEY,required"`
}

func (x Values) Name(v ...string) string {
	return fmt.Sprintf(`%s:%s`, x.Namespace, strings.Join(v, ":"))
}

func GetClaims(c *app.RequestContext) (claims passport.Claims) {
	value, ok := c.Get("identity")
	if !ok {
		return
	}
	return value.(passport.Claims)
}
