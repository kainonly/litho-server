package common

import (
	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/nats-io/nats.go"
	"github.com/tencentyun/cos-go-sdk-v5"
	"github.com/weplanx/go/encryption"
	"github.com/weplanx/go/passport"
	"go.mongodb.org/mongo-driver/mongo"
)

type Inject struct {
	Values      *Values
	MongoClient *mongo.Client
	Db          *mongo.Database
	Redis       *redis.Client
	Nats        *nats.Conn
	Js          nats.JetStreamContext
	Passport    *passport.Passport
	Cipher      *encryption.Cipher
	HID         *encryption.HID
	Cos         *cos.Client
}

const TokenClaimsKey = "token-claims"

var (
	LoginExpired = errors.New("login token has expired")
)
