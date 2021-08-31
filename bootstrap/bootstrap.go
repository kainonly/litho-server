package bootstrap

import (
	"context"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/kainonly/go-bit/authx"
	"github.com/kainonly/go-bit/cipher"
	"github.com/kainonly/go-bit/cookie"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"lab-api/common"
	"net/http"
	"time"
)

// HttpServer 启动 Gin HTTP 服务
// 配置文档 https://gin-gonic.com/docs/examples/custom-http-config
func HttpServer(config common.Config) (serve *gin.Engine) {
	serve = gin.New()
	serve.Use(gin.Logger())
	serve.Use(gin.Recovery())
	serve.Use(cors.New(cors.Config{
		AllowOrigins:     config.Cors,
		AllowMethods:     []string{"POST"},
		AllowHeaders:     []string{"Origin", "CONTENT-TYPE"},
		AllowCredentials: true,
		MaxAge:           6 * time.Hour,
	}))
	return
}

// InitializeDatabase 初始化数据库
// 配置文档 https://gorm.io/docs/connecting_to_the_database.html
func InitializeDatabase(config common.Config) (db *gorm.DB, err error) {
	option := config.Database
	db, err = gorm.Open(postgres.Open(option.Dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
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
func InitializeRedis(config common.Config) (client *redis.Client, err error) {
	option := config.Redis
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
func InitializeCookie(config common.Config) *cookie.Cookie {
	return cookie.New(config.Cookie, http.SameSiteStrictMode)
}

// InitializeAuthx 创建认证
func InitializeAuthx(config common.Config) *authx.Authx {
	options := map[string]*authx.Auth{
		"system": {
			Key: config.App.Key,
			Iss: config.App.Name,
			Aud: []string{"admin"},
			Exp: 720,
		},
	}
	return authx.New(options)
}

// InitializeCipher 初始化数据加密
func InitializeCipher(config common.Config) (*cipher.Cipher, error) {
	return cipher.New(config.App.Key)
}
