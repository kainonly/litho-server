package common

import (
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/redis/go-redis/v9"
	"github.com/weplanx/go-wpx/passport"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/fx"
	"strings"
)

type Inject struct {
	fx.In

	V     *Values
	Hertz *server.Hertz
	Mgo   *mongo.Client
	Db    *mongo.Database
	RDb   *redis.Client
}

type Values struct {
	Address   string `env:"ADDRESS" envDefault:":3000"`
	Namespace string `env:"NAMESPACE,required"`
	Key       string `env:"KEY,required"`
	Database  `envPrefix:"DATABASE_"`
	Nats      `envPrefix:"NATS_"`
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
