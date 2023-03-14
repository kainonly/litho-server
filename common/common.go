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
	Values    *Values
	Mongo     *mongo.Client
	Db        *mongo.Database
	Redis     *redis.Client
	Influx    influxdb2.Client
	Nats      *nats.Conn
	JetStream nats.JetStreamContext
	KeyValue  nats.KeyValue
}

func GetClaims(c *app.RequestContext) (claims passport.Claims) {
	value, ok := c.Get("identity")
	if !ok {
		return
	}
	return value.(passport.Claims)
}
