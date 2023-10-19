package common

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"github.com/bytedance/sonic"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/imroc/req/v3"
	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"
	transfer "github.com/weplanx/collector/client"
	"github.com/weplanx/go/captcha"
	"github.com/weplanx/go/cipher"
	"github.com/weplanx/go/locker"
	"github.com/weplanx/go/passport"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type Inject struct {
	V         *Values
	Mgo       *mongo.Client
	Db        *mongo.Database
	RDb       *redis.Client
	Nats      *nats.Conn
	JetStream nats.JetStreamContext
	KeyValue  nats.KeyValue
	Cipher    *cipher.Cipher
	Captcha   *captcha.Captcha
	Locker    *locker.Locker
	Transfer  *transfer.Client
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

func HttpClient(url string) *req.Client {
	return req.C().
		SetBaseURL(url).
		SetJsonMarshal(sonic.Marshal).
		SetJsonUnmarshal(sonic.Unmarshal).
		SetTimeout(time.Second * 5)
}

func Sha256hex(s string) string {
	b := sha256.Sum256([]byte(s))
	return hex.EncodeToString(b[:])
}

func Hmacsha256(s, key string) string {
	hashed := hmac.New(sha256.New, []byte(key))
	hashed.Write([]byte(s))
	return string(hashed.Sum(nil))
}
