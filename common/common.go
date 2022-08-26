package common

import (
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/go-redis/redis/v8"
	"github.com/nats-io/nats.go"
	"github.com/weplanx/support"
	"go.mongodb.org/mongo-driver/mongo"
)

type Inject struct {
	Values *support.Values
	Mongo  *mongo.Client
	Db     *mongo.Database
	Redis  *redis.Client
	Nats   *nats.Conn
}

// Active 授权用户标识
type Active struct {
	// Token ID
	JTI string

	// User ID
	UID string
}

// GetActive 获取授权用户标识
func GetActive(c *app.RequestContext) (data Active) {
	value, ok := c.Get("identity")
	if !ok {
		return
	}
	return value.(Active)
}
