package bootstrap

import (
	"api/common"
	"context"
	"errors"
	"fmt"
	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
	"github.com/speps/go-hashids/v2"
	"github.com/tencentyun/cos-go-sdk-v5"
	"github.com/weplanx/go/encryption"
	"github.com/weplanx/go/engine"
	"github.com/weplanx/go/passport"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"
)

var Provides = wire.NewSet(
	UseMongoDB,
	UseDatabase,
	UseRedis,
	UsePulsar,
	UseEngine,
	UsePassport,
	UseCipher,
	UseIDx,
	UseCos,
)

// SetValues 初始化配置
func SetValues() (values *common.Values, err error) {
	if _, err = os.Stat("./config/config.yml"); os.IsNotExist(err) {
		return nil, errors.New("静态配置不存在，请检查路径 [./config/config.yml]")
	}
	var config []byte
	if config, err = ioutil.ReadFile("./config/config.yml"); err != nil {
		return
	}
	if err = yaml.Unmarshal(config, &values); err != nil {
		return
	}
	return
}

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

// UsePulsar 初始化 Pulsar
func UsePulsar(values *common.Values) (pulsar.Client, error) {
	option := values.Pulsar
	return pulsar.NewClient(pulsar.ClientOptions{
		URL:               option.Url,
		Authentication:    pulsar.NewAuthenticationToken(option.Token),
		OperationTimeout:  30 * time.Second,
		ConnectionTimeout: 30 * time.Second,
	})
}

// UseEngine 初始化 Weplanx Engine
func UseEngine(values *common.Values, client pulsar.Client) *engine.Engine {
	return engine.New(
		engine.SetApp(values.Name),
		engine.UseStaticOptions(values.Engines),
		engine.UsePulsar(client),
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
	values.Passport.Iss = values.Name
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
