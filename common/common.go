package common

import (
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
)

type Inject struct {
	Values *Values
	Mongo  *mongo.Client
	Db     *mongo.Database
	Redis  *redis.Client
}
