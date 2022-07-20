package app

import (
	"github.com/go-redis/redis/v8"
	"github.com/weplanx/server/common"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type Service struct {
	Values *common.Values
	Mongo  *mongo.Client
	Db     *mongo.Database
	Redis  *redis.Client
}

func (x *Service) Index() time.Time {
	return time.Now()
}
