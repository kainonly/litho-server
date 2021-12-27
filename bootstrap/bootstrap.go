package bootstrap

import (
	"api/common"
	"context"
	"errors"
	"github.com/google/wire"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
)

var Provides = wire.NewSet(
	UseMongoDB,
	UseDatabase,
)

// SetValues 初始化配置
func SetValues() (values *common.Values, err error) {
	if _, err = os.Stat("./config.yml"); os.IsNotExist(err) {
		err = errors.New("the path [./config.yml] does not have a configuration file")
		return
	}
	var b []byte
	b, err = ioutil.ReadFile("./config.yml")
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

//// UseRedis 初始化Redis缓存
//// 配置文档 https://github.com/go-redis/redis
//func UseRedis(values *common.Values) (client *redis.Client, err error) {
//	opts, err := redis.ParseURL(values.Redis.Uri)
//	if err != nil {
//		return
//	}
//	client = redis.NewClient(opts)
//	if err = client.Ping(context.Background()).Err(); err != nil {
//		return
//	}
//	return
//}

//// UsePassport 创建认证
//func UsePassport(values *common.Values) *passport.Passport {
//	values.Passport.Iss = values.Name
//	return passport.New(values.Key, values.Passport)
//}

//func UseEncryption(values *common.Values) (cipher *encryption.Cipher, idx *encryption.IDx, err error) {
//	if cipher, err = encryption.NewCipher(values.Key); err != nil {
//		return
//	}
//	if idx, err = encryption.NewIDx(values.Key); err != nil {
//		return
//	}
//	return
//}
