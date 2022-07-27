package common

import (
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/go-redis/redis/v8"
	"github.com/nats-io/nats.go"
	"go.mongodb.org/mongo-driver/mongo"
)

type Inject struct {
	Values *Values
	Mongo  *mongo.Client
	Db     *mongo.Database
	Redis  *redis.Client
	Nats   *nats.Conn
}

type Active struct {
	JTI string
	UID string
}

func GetActive(c *app.RequestContext) (data Active) {
	value, ok := c.Get("identity")
	if !ok {
		return
	}
	return value.(Active)
}
