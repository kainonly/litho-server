package common

import (
	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/nats-io/nats.go"
	"github.com/tencentyun/cos-go-sdk-v5"
	"github.com/weplanx/go/encryption"
	"github.com/weplanx/go/passport"
	openapi "github.com/weplanx/openapi/client"
	"github.com/weplanx/transfer"
	"go.mongodb.org/mongo-driver/mongo"
)

type Inject struct {
	Values      *Values
	MongoClient *mongo.Client
	Db          *mongo.Database
	Redis       *redis.Client
	Nats        *nats.Conn
	Js          nats.JetStreamContext
	Transfer    *transfer.Transfer
	Open        *openapi.OpenAPI
	Passport    *passport.Passport
	Cipher      *encryption.Cipher
	HID         *encryption.HID
	Cos         *cos.Client
}

const TokenClaimsKey = "token-claims"

var (
	AuthExpired  = errors.New("认证已失效，令牌超出有效期")
	AuthConflict = errors.New("认证已失效，已被新终端占用")
)
