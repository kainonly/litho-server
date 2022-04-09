package bootstrap

import (
	"api/common"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nkeys"
	"github.com/speps/go-hashids/v2"
	"github.com/tencentyun/cos-go-sdk-v5"
	"github.com/weplanx/go/encryption"
	"github.com/weplanx/go/engine"
	"github.com/weplanx/go/passport"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var Provides = wire.NewSet(
	UseMongoDB,
	UseDatabase,
	UseRedis,
	UseNats,
	UseJetStream,
	UseEngine,
	UsePassport,
	UseCipher,
	UseIDx,
	UseCos,
)

// UseMongoDB 初始化 Mongodb
func UseMongoDB(values *common.Values) (*mongo.Client, error) {
	return mongo.Connect(
		context.TODO(),
		options.Client().ApplyURI(values.Database.Uri),
	)
}

// UseDatabase 指向使用数据库名称
func UseDatabase(client *mongo.Client, values *common.Values) (db *mongo.Database) {
	return client.Database(values.Database.DbName)
}

// UseRedis 初始化 Redis 缓存
// 配置文档 https://github.com/go-redis/redis
func UseRedis(values *common.Values) (client *redis.Client, err error) {
	opts, err := redis.ParseURL(values.Redis.Uri)
	if err != nil {
		return
	}
	client = redis.NewClient(opts)
	if err = client.Ping(context.TODO()).Err(); err != nil {
		return
	}
	return
}

func UseNats(values *common.Values) (nc *nats.Conn, err error) {
	var kp nkeys.KeyPair
	if kp, err = nkeys.FromSeed([]byte(values.Nats.Nkey)); err != nil {
		return
	}
	defer kp.Wipe()
	var pub string
	if pub, err = kp.PublicKey(); err != nil {
		return
	}
	if !nkeys.IsValidPublicUserKey(pub) {
		return nil, fmt.Errorf("nkey 验证失败")
	}
	if nc, err = nats.Connect(
		strings.Join(values.Nats.Hosts, ","),
		nats.MaxReconnects(5),
		nats.ReconnectWait(2*time.Second),
		nats.ReconnectJitter(500*time.Millisecond, 2*time.Second),
		nats.Nkey(pub, func(nonce []byte) ([]byte, error) {
			sig, _ := kp.Sign(nonce)
			return sig, nil
		}),
	); err != nil {
		return
	}
	return
}

func UseJetStream(nc *nats.Conn) (nats.JetStreamContext, error) {
	return nc.JetStream(nats.PublishAsyncMaxPending(256))
}

// UseEngine 初始化 Weplanx Engine
func UseEngine(values *common.Values, js nats.JetStreamContext) *engine.Engine {
	return engine.New(
		engine.SetApp(values.Namespace),
		engine.UseStaticOptions(values.Engines),
		engine.UseEvents(js),
	)
}

// UseCipher 使用数据加密
func UseCipher(values *common.Values) (cipher *encryption.Cipher, err error) {
	if cipher, err = encryption.NewCipher(values.Key); err != nil {
		return
	}
	return
}

// UseIDx 使用 ID 加密
func UseIDx(values *common.Values) (idx *encryption.IDx, err error) {
	if idx, err = encryption.NewIDx(values.Key, hashids.DefaultAlphabet); err != nil {
		return
	}
	return
}

// UsePassport 创建认证
func UsePassport(values *common.Values) *passport.Passport {
	values.Passport.Iss = values.Namespace
	return passport.New(values.Key, values.Passport)
}

func UseCos(values *common.Values) (client *cos.Client, err error) {
	option := values.QCloud
	var u *url.URL
	if u, err = url.Parse(
		fmt.Sprintf(`https://%s.cos.%s.myqcloud.com`, option.Cos.Bucket, option.Cos.Region),
	); err != nil {
		return
	}
	client = cos.NewClient(&cos.BaseURL{BucketURL: u}, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  option.SecretID,
			SecretKey: option.SecretKey,
		},
	})
	return
}
