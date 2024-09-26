package common

import (
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"
	"github.com/weplanx/go/captcha"
	"github.com/weplanx/go/cipher"
	"github.com/weplanx/go/locker"
	"github.com/weplanx/go/passport"
)

type Inject struct {
	V         *Values
	RDb       *redis.Client
	Nats      *nats.Conn
	JetStream nats.JetStreamContext
	KeyValue  nats.KeyValue
	Cipher    *cipher.Cipher
	Captcha   *captcha.Captcha
	Locker    *locker.Locker
}

type APIPassport = passport.Passport

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
