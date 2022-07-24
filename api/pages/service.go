package pages

import (
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
)

type Service struct {
	Db    *mongo.Database
	Redis *redis.Client
}
