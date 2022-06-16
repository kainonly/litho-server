package common

import (
	"github.com/go-redis/redis/v8"
	"github.com/go-resty/resty/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"sync"
)

type Inject struct {
	Values      *Values
	MongoClient *mongo.Client
	Db          *mongo.Database
	Redis       *redis.Client
}

type Jobs = sync.Map

type HttpClients struct {
	Feishu *resty.Client
}
