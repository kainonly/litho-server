package bootstrap

import (
	"context"
	"errors"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/kainonly/go-bit/authx"
	"github.com/kainonly/go-bit/cipher"
	"github.com/kainonly/go-bit/cookie"
	"github.com/kainonly/go-bit/crud"
	"go.uber.org/fx"
	"gopkg.in/yaml.v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"io/ioutil"
	"lab-api/common"
	"net/http"
	"os"
	"time"
)

var Provides = fx.Provide(
	LoadSettings,
	InitializeDatabase,
	crud.New,
	InitializeRedis,
	InitializeCookie,
	InitializeAuthx,
	InitializeCipher,
	HttpServer,
)

// LoadSettings 初始化应用配置
func LoadSettings() (app *common.Set, err error) {
	if _, err = os.Stat("./config.yml"); os.IsNotExist(err) {
		err = errors.New("当前路径 [./config.yml] 不存在配置文件")
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

// InitializeDatabase 初始化数据库
// 配置文档 https://gorm.io/docs/connecting_to_the_database.html
func InitializeDatabase(app *common.Set) (db *gorm.DB, err error) {
	option := app.Database
	db, err = gorm.Open(postgres.Open(option.Dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		return
	}
	sqlDB, err := db.DB()
	if err != nil {
		return
	}
	if option.MaxIdleConns != 0 {
		sqlDB.SetMaxIdleConns(option.MaxIdleConns)
	}
	if option.MaxOpenConns != 0 {
		sqlDB.SetMaxOpenConns(option.MaxOpenConns)
	}
	if option.ConnMaxLifetime != 0 {
		sqlDB.SetConnMaxLifetime(time.Second * time.Duration(option.ConnMaxLifetime))
	}
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
	if err = client.Ping(context.Background()).Err(); err != nil {
		return
	}
	return
}

// InitializeCookie 创建 Cookie 工具
func InitializeCookie(app *common.Set) *cookie.Cookie {
	return cookie.New(app.Cookie, http.SameSiteStrictMode)
}

// InitializeAuthx 创建认证
func InitializeAuthx(app *common.Set) *authx.Authx {
	options := map[string]*authx.Auth{
		"system": {
			Key: app.Key,
			Iss: app.Name,
			Aud: []string{"admin"},
			Exp: 720,
		},
	}
	return authx.New(options)
}

// InitializeCipher 初始化数据加密
func InitializeCipher(app *common.Set) (*cipher.Cipher, error) {
	return cipher.New(app.Key)
}

// HttpServer 启动 Gin HTTP 服务
// 配置文档 https://gin-gonic.com/docs/examples/custom-http-config
func HttpServer(lc fx.Lifecycle, config *common.Set) (router *gin.Engine) {
	router = gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(cors.New(cors.Config{
		AllowOrigins:     config.Cors,
		AllowMethods:     []string{"POST"},
		AllowHeaders:     []string{"Origin", "CONTENT-TYPE"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go router.Run(":9000")
			return nil
		},
	})
	return
}
