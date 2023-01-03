package common

import (
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/go-redis/redis/v8"
	"github.com/go-resty/resty/v2"
	"github.com/nats-io/nats.go"
	"github.com/weplanx/utils/kv"
	"github.com/weplanx/utils/passport"
	"go.mongodb.org/mongo-driver/mongo"
	"strings"
)

type Inject struct {
	Values      *Values
	Mongo       *mongo.Client
	Db          *mongo.Database
	Redis       *redis.Client
	Nats        *nats.Conn
	JetStream   nats.JetStreamContext
	KeyValue    nats.KeyValue
	HttpClients *HttpClients
}

type Values struct {
	Address           string   `env:"ADDRESS" envDefault:":3000"`
	Namespace         string   `env:"NAMESPACE,required"`
	Key               string   `env:"KEY,required"`
	Console           string   `env:"CONSOLE,required"`
	Hosts             []string `env:"HOSTS" envSeparator:","`
	Database          `envPrefix:"DATABASE_"`
	Nats              `envPrefix:"NATS_"`
	*kv.DynamicValues `env:"-"`
}

type Database struct {
	Mongo string `env:"MONGO,required"`
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

type HttpClients struct {
	Feishu *resty.Client
}

func GetClaims(c *app.RequestContext) (claims passport.Claims) {
	value, ok := c.Get("identity")
	if !ok {
		return
	}
	return value.(passport.Claims)
}
