package common

import (
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"
	transfer "github.com/weplanx/collector/client"
	"github.com/weplanx/go/captcha"
	"github.com/weplanx/go/cipher"
	"github.com/weplanx/go/locker"
	"github.com/weplanx/go/passport"
	"go.mongodb.org/mongo-driver/mongo"
)

type Inject struct {
	V         *Values
	Mgo       *mongo.Client
	Db        *mongo.Database
	RDb       *redis.Client
	Flux      influxdb2.Client
	Nats      *nats.Conn
	JetStream nats.JetStreamContext
	KeyValue  nats.KeyValue
	Cipher    *cipher.Cipher
	Passport  *passport.Passport
	Captcha   *captcha.Captcha
	Locker    *locker.Locker
	Transfer  *transfer.Client
}

func Claims(c *app.RequestContext) (claims passport.Claims) {
	value, ok := c.Get("identity")
	if !ok {
		return
	}
	return value.(passport.Claims)
}

func SetAccessToken(c *app.RequestContext, ts string) {
	c.SetCookie("TOKEN", ts, -1,
		"/", "", protocol.CookieSameSiteLaxMode, true, true)
}

func ClearAccessToken(c *app.RequestContext) {
	c.SetCookie("TOKEN", "", -1,
		"/", "", protocol.CookieSameSiteLaxMode, true, true)
}
