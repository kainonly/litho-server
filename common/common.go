package common

import (
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/go-redis/redis/v8"
	"github.com/go-resty/resty/v2"
	"github.com/nats-io/nats.go"
	"github.com/weplanx/utils/passport"
	"go.mongodb.org/mongo-driver/mongo"
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
