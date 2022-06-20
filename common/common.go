package common

import (
	"github.com/go-redis/redis/v8"
	"github.com/go-resty/resty/v2"
	"github.com/nats-io/nats.go"
	"go.mongodb.org/mongo-driver/mongo"
	"sync"
)

type Inject struct {
	Values      *Values
	MongoClient *mongo.Client
	Db          *mongo.Database
	Redis       *redis.Client
	Store       nats.ObjectStore
	HC          *HttpClients
}

type Jobs = sync.Map

type HttpClients struct {
	Feishu *resty.Client
}
