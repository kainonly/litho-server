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
		err = errors.New("the path [./config.yml] does not have a configuration file")
		return
	}
	var b []byte
	b, err = ioutil.ReadFile("./config/config.yml")
	if err != nil {
		return
	}
	err = yaml.Unmarshal(b, &values)
	if err != nil {
		return
	}
	return
}

func UseMongoDB(values *common.Values) (*mongo.Client, error) {
	return mongo.Connect(
		context.TODO(),
		options.Client().ApplyURI(values.Database.Uri),
	)
}

func UseDatabase(client *mongo.Client, values *common.Values) (db *mongo.Database) {
	return client.Database(values.Database.DbName)
}

// UseRedis 初始化Redis缓存
// 配置文档 https://github.com/go-redis/redis
func UseRedis(values *common.Values) (client *redis.Client, err error) {
	opts, err := redis.ParseURL(values.Redis.Uri)
	if err != nil {
		return
	}
	client = redis.NewClient(opts)
	if err = client.Ping(context.Background()).Err(); err != nil {
		return
	}
	return
}

func UsePulsar(values *common.Values) (pulsar.Client, error) {
	option := values.Pulsar
	return pulsar.NewClient(pulsar.ClientOptions{
		URL:               option.Url,
		Authentication:    pulsar.NewAuthenticationToken(option.Token),
		OperationTimeout:  30 * time.Second,
		ConnectionTimeout: 30 * time.Second,
	})
}

func UseEngine(values *common.Values, client pulsar.Client) *engine.Engine {
	return engine.New(
		engine.SetApp(values.Name),
		engine.UseStaticOptions(values.Engines),
		engine.UsePulsar(client),
	)
}

// UsePassport 创建认证
func UsePassport(values *common.Values) *passport.Passport {
	values.Passport.Iss = values.Name
	return passport.New(values.Key, values.Passport)
}

func UseCipher(values *common.Values) (cipher *encryption.Cipher, err error) {
	if cipher, err = encryption.NewCipher(values.Key); err != nil {
		return
	}
	return
}

func UseIDx(values *common.Values) (idx *encryption.IDx, err error) {
	if idx, err = encryption.NewIDx(values.Key, hashids.DefaultAlphabet); err != nil {
		return
	}
	return
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
