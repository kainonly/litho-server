package common

import (
	"github.com/go-redis/redis/v8"
	"github.com/go-resty/resty/v2"
	"github.com/nats-io/nats.go"
	"github.com/weplanx/go/encryption"
	"github.com/weplanx/go/passport"
	"go.mongodb.org/mongo-driver/mongo"
	"sync"
)

type Inject struct {
	Values      *Values
	MongoClient *mongo.Client
	Db          *mongo.Database
	Redis       *redis.Client
	Nats        *nats.Conn
	Js          nats.JetStreamContext
	Cipher      *encryption.Cipher
	HID         *encryption.HID
	Passport    *passport.Passport
	HC          *HttpClients
}

type HttpClients struct {
	Feishu *resty.Client
}

type Subscriptions struct {
	*sync.Map
}
