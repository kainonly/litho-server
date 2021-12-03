package bootstrap

import (
	"api/common"
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/encryptcookie"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/weplanx/go/api"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/fx"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

var Provides = fx.Provide(
	SetValues,
	UseDatabase,
	UseRedis,
	HttpServer,
	api.New,
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

//// InitializePassport 创建认证
//func UsePassport(values *common.Values) *passport.Passport {
//	return passport.New(map[string]*passport.Auth{
//		"system": {
//			Key: values.Key,
//			Iss: values.Name,
//			Aud: []string{"admin"},
//			Exp: 720,
//		},
//	})
//}

// HttpServer 启动 HTTP 服务
func HttpServer(lc fx.Lifecycle, values *common.Values) (app *fiber.App) {
	app = fiber.New(fiber.Config{
		AppName: values.Name,
	})
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     strings.Join(values.Cors.AllowOrigins, ","),
		AllowMethods:     strings.Join(values.Cors.AllowMethods, ","),
		AllowHeaders:     strings.Join(values.Cors.AllowHeaders, ","),
		AllowCredentials: values.Cors.AllowCredentials,
		MaxAge:           int(time.Duration(values.Cors.MaxAge) * time.Second),
	}))
	app.Use(encryptcookie.New(encryptcookie.Config{
		Key: values.Key,
	}))
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go app.Listen(values.Address)
			return nil
		},
	})
	return
}
