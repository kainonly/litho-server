package bootstrap

import (
	"api/common"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/go-resty/resty/v2"
	"github.com/google/wire"
	jsoniter "github.com/json-iterator/go"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nkeys"
	"github.com/speps/go-hashids/v2"
	"github.com/weplanx/go/encryption"
	"github.com/weplanx/go/engine"
	"github.com/weplanx/go/passport"
	"github.com/weplanx/transfer"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strings"
	"time"
)

var Provides = wire.NewSet(
	UseMongoDB,
	UseDatabase,
	UseRedis,
	UseNats,
	UseJetStream,
	UseStore,
	UseEngine,
	UseTransfer,
	UsePassport,
	UseCipher,
	UseHID,
	UseHttpClients,
)

// UseMongoDB 初始化 MongoDB
// 配置文档 https://www.mongodb.com/docs/drivers/go/current/
func UseMongoDB(values *common.Values) (*mongo.Client, error) {
	return mongo.Connect(
		context.TODO(),
		options.Client().ApplyURI(values.Database.Uri),
	)
}

// UseDatabase 初始化数据库
func UseDatabase(client *mongo.Client, values *common.Values) (db *mongo.Database) {
	return client.Database(values.Database.DbName)
}

// UseRedis 初始化 Redis
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

// UseNats 初始化 Nats
// 配置文档 https://docs.nats.io/using-nats/developer
// SDK https://github.com/nats-io/nats.go
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

// UseJetStream 初始化流
func UseJetStream(nc *nats.Conn) (nats.JetStreamContext, error) {
	return nc.JetStream(nats.PublishAsyncMaxPending(256))
}

// UseStore 初始化分布存储
func UseStore(values *common.Values, js nats.JetStreamContext) (store nats.ObjectStore, err error) {
	if store, err = js.CreateObjectStore(&nats.ObjectStoreConfig{
		Bucket: values.Namespace,
	}); err != nil {
		return
	}
	// 初始化动态配置
	var b []byte
	if b, err = store.GetBytes("values"); err != nil {
		if err == nats.ErrObjectNotFound {
			dv := &common.DynamicValues{
				UserSessionExpire:    time.Hour,
				UserLoginFailedTimes: 5,
				UserLockTime:         time.Minute * 15,
				IpLoginFailedTimes:   10,
				IpWhitelist:          []string{},
				IpBlacklist:          []string{},
				PasswordStrength:     1,
				PasswordExpire:       365,
				TencentCosExpired:    time.Second * 300,
				TencentCosLimit:      5120,
				EmailPort:            "465",
			}
			if b, err = jsoniter.Marshal(dv); err != nil {
				return
			}
			if _, err = store.PutBytes("values", b); err != nil {
				return
			}
			return
		} else {
			return
		}
	}
	if err = jsoniter.Unmarshal(b, &values.DynamicValues); err != nil {
		return
	}
	// 监听配置
	go func() {
		var watch nats.ObjectWatcher
		if watch, err = store.Watch(); err != nil {
			return
		}
		current := time.Now()
		for o := range watch.Updates() {
			if o == nil || o.ModTime.Unix() < current.Unix() {
				continue
			}
			if o.Name == "values" {
				var b []byte
				b, err = store.GetBytes("values")
				if err != nil {
					return
				}
				if err = jsoniter.Unmarshal(b, &values.DynamicValues); err != nil {
					// TODO: 发送异常提示
					return
				}
			}
		}
	}()
	return
}

// UseEngine 初始化 Engine API
func UseEngine(values *common.Values, js nats.JetStreamContext) *engine.Engine {
	return engine.New(
		engine.SetApp(values.Namespace),
		engine.UseStaticOptions(values.Engines),
		engine.UseEvents(js),
	)
}

// UseTransfer 初始化日志传输
func UseTransfer(values *common.Values, js nats.JetStreamContext) (client *transfer.Transfer, err error) {
	if client, err = transfer.New(values.Namespace, js); err != nil {
		return nil, err
	}
	if err = client.Set("request", transfer.Option{
		Topic:       "request",
		Description: "请求日志",
	}); err != nil {
		return
	}
	return
}

// UsePassport 创建认证
func UsePassport(values *common.Values) *passport.Passport {
	return passport.New(values.Key, passport.Option{
		Iss: values.Namespace,
		Aud: []string{"console"},
		Exp: 720,
	})
}

// UseCipher 数据加密
func UseCipher(values *common.Values) (cipher *encryption.Cipher, err error) {
	if cipher, err = encryption.NewCipher(values.Key); err != nil {
		return
	}
	return
}

// UseHID ID加密
func UseHID(values *common.Values) (idx *encryption.HID, err error) {
	if idx, err = encryption.NewIDx(values.Key, hashids.DefaultAlphabet); err != nil {
		return
	}
	return
}

// UseHttpClients 创建请求客户端
func UseHttpClients() *common.HttpClients {
	return &common.HttpClients{
		Feishu: resty.New().SetBaseURL("https://open.feishu.cn/open-apis"),
	}
}
