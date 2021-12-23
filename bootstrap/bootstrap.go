package bootstrap

import (
	"api/common"
	"context"
	"errors"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/weplanx/go/api"
	"github.com/weplanx/go/encryption"
	"github.com/weplanx/go/passport"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/fx"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"time"
)

var Provides = fx.Provide(
	SetValues,
	UseDatabase,
	UseRedis,
	UsePassport,
	UseEncryption,
	HttpServer,
	api.New,
	api.AutoController,
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

func UseDatabase(values *common.Values) (client *mongo.Client, db *mongo.Database, err error) {
	if client, err = mongo.Connect(
		context.TODO(),
		options.Client().ApplyURI(values.Database.Uri),
	); err != nil {
		return
	}
	db = client.Database(values.Database.DbName)
	return
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

// UsePassport 创建认证
func UsePassport(values *common.Values) *passport.Passport {
	values.Passport.Iss = values.Name
	return passport.New(values.Key, values.Passport)
}

func UseEncryption(values *common.Values) (cipher *encryption.Cipher, idx *encryption.IDx, err error) {
	if cipher, err = encryption.NewCipher(values.Key); err != nil {
		return
	}
	if idx, err = encryption.NewIDx(values.Key); err != nil {
		return
	}
	return
}

// HttpServer 启动 HTTP 服务
func HttpServer(lc fx.Lifecycle, values *common.Values) (r *gin.Engine) {
	r = gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     values.Cors.AllowOrigins,
		AllowMethods:     values.Cors.AllowMethods,
		AllowHeaders:     values.Cors.AllowHeaders,
		AllowCredentials: values.Cors.AllowCredentials,
		MaxAge:           time.Duration(values.Cors.MaxAge) * time.Second,
	}))
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go r.Run(values.Address)
			return nil
		},
	})
	return
}
