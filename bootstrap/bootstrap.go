package bootstrap

import (
	"api/common"
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
	"github.com/weplanx/go/api"
	"github.com/weplanx/go/helper"
	"github.com/weplanx/go/passport"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"net/http"
	"os"
)

var Provides = wire.NewSet(
	LoadSettings,
	InitializeMongoDB,
	InitializeDatabase,
	InitializeRedis,
	InitializeCommonApi,
	InitializeCookie,
	InitializePassport,
	InitializeCipher,
)

// LoadSettings 初始化应用配置
func LoadSettings() (app *common.Set, err error) {
	if _, err = os.Stat("./config.yml"); os.IsNotExist(err) {
		err = errors.New("the path [./config.yml] does not have a configuration file")
		return
	}
	var b []byte
	b, err = ioutil.ReadFile("./config.yml")
	if err != nil {
		return
	}
	err = yaml.Unmarshal(b, &app)
	if err != nil {
		return
	}
	return
}

func InitializeMongoDB(app *common.Set) (client *mongo.Client, err error) {
	option := app.Database
	if client, err = mongo.Connect(
		context.TODO(),
		options.Client().ApplyURI(option.Uri),
	); err != nil {
		return
	}
	return
}

// InitializeDatabase 初始化数据库
func InitializeDatabase(app *common.Set, client *mongo.Client) (db *mongo.Database) {
	option := app.Database
	db = client.Database(option.Name)
	return
}

// InitializeRedis 初始化Redis缓存
// 配置文档 https://github.com/go-redis/redis
func InitializeRedis(app *common.Set) (client *redis.Client, err error) {
	option := app.Redis
	client = redis.NewClient(&redis.Options{
		Addr:     option.Address,
		Password: option.Password,
		DB:       option.DB,
	})
	// Serverless 模式建议关闭
	//if err = client.Ping(context.Background()).Err(); err != nil {
	//	return
	//}
	return
}

func InitializeCommonApi(client *mongo.Client, db *mongo.Database) *api.API {
	return api.New(client, db)
}

// InitializeCookie 创建 Cookie 工具
func InitializeCookie(app *common.Set) *helper.CookieHelper {
	return helper.NewCookieHelper(app.Cookie, http.SameSiteStrictMode)
}

// InitializePassport 创建认证
func InitializePassport(app *common.Set) *passport.Passport {
	return passport.New(map[string]*passport.Auth{
		"system": {
			Key: app.Key,
			Iss: app.Name,
			Aud: []string{"admin"},
			Exp: 720,
		},
	})
}

// InitializeCipher 初始化数据加密
func InitializeCipher(app *common.Set) (*helper.CipherHelper, error) {
	return helper.NewCipherHelper(app.Key)
}
