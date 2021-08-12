package bootstrap

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/kainonly/go-bit"
	"github.com/mitchellh/mapstructure"
	"go.uber.org/fx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"time"
)

// InitializeDatabase 初始化 Postgresql 数据库
// 配置文档 https://gorm.io/docs/connecting_to_the_database.html
func InitializeDatabase(config bit.Config) (db *gorm.DB, err error) {
	var option struct {
		Dsn             string
		MaxIdleConns    int
		MaxOpenConns    int
		ConnMaxLifetime int
		TablePrefix     string
	}
	if err = mapstructure.Decode(config["database"], &option); err != nil {
		return
	}
	db, err = gorm.Open(postgres.Open(option.Dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   option.TablePrefix,
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
func InitializeRedis(config bit.Config) (client *redis.Client, err error) {
	var option struct {
		Address  string
		Password string
		DB       int
	}
	if err = mapstructure.Decode(config["redis"], &option); err != nil {
		return
	}
	client = redis.NewClient(&redis.Options{
		Addr:     option.Address,
		Password: option.Password,
		DB:       option.DB,
	})
	return
}

// HttpServer 启动 Gin HTTP 服务
// 配置文档 https://gin-gonic.com/docs/examples/custom-http-config
func HttpServer(lc fx.Lifecycle) (serve *gin.Engine) {
	serve = gin.New()
	serve.Use(gin.Logger())
	serve.Use(gin.Recovery())
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go serve.Run(":8000")
			return nil
		},
	})
	return
}
