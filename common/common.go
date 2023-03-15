package common

import (
	"github.com/cloudwego/hertz/pkg/app"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"
	"github.com/weplanx/utils/passport"
	"go.mongodb.org/mongo-driver/mongo"
)

type Inject struct {
	V    *Values
	Mgo  *mongo.Client
	Db   *mongo.Database
	RDb  *redis.Client
	Flux influxdb2.Client
	Nats *nats.Conn
	JS   nats.JetStreamContext
	KV   nats.KeyValue
}

func GetClaims(c *app.RequestContext) (claims passport.Claims) {
	value, ok := c.Get("identity")
	if !ok {
		return
	}
	return value.(passport.Claims)
}
