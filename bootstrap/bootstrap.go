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
	"gopkg.in/yaml.v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"io/ioutil"
	"lab-api/config"
	"os"
	"time"
)

// LoadConfiguration 初始化应用配置
func LoadConfiguration() (config config.Config, err error) {
	if _, err = os.Stat("./config.yml"); os.IsNotExist(err) {
		err = errors.New("the configuration file path [./config.yml] does not exist")
		return
	}
	var buf []byte
	buf, err = ioutil.ReadFile("./config.yml")
	if err != nil {
		return
	}
	err = yaml.Unmarshal(buf, &config)
	if err != nil {
		return
	}
	return
}

// InitializeDatabase 初始化 Postgresql 数据库
// 配置文档 https://gorm.io/docs/connecting_to_the_database.html
func InitializeDatabase(config config.Config) (db *gorm.DB, err error) {
	option := config.Database
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
func InitializeRedis(config config.Config) (client *redis.Client, err error) {
	option := config.Redis
	client = redis.NewClient(&redis.Options{
		Addr:     option.Address,
		Password: option.Password,
		DB:       option.DB,
	})
	return
}

// InitializeCrud 初始化 CRUD 工具
func InitializeCrud(db *gorm.DB) *crud.Crud {
	return crud.New(db)
}

// InitializeCipher 初始化数据加密
func InitializeCipher(config config.Config) (*cipher.Cipher, error) {
	return cipher.New(config.App.Key)
}

// InitializeCookie 创建 Cookie 工具
func InitializeCookie(config config.Config) *cookie.Cookie {
	return cookie.New(config.Cookie)
}

func InitializeAuthx(config config.Config) *authx.Authx {
	options := config.Auth
	for _, v := range options {
		if v.Key == "" {
			v.Key = config.App.Key
		}
		if v.Iss == "" {
			v.Iss = config.App.Name
		}
	}
	return authx.New(options)
}

// HttpServer 启动 Gin HTTP 服务
// 配置文档 https://gin-gonic.com/docs/examples/custom-http-config
func HttpServer(lc fx.Lifecycle, config config.Config) (serve *gin.Engine) {
	serve = gin.New()
	serve.Use(gin.Logger())
	serve.Use(gin.Recovery())
	serve.Use(cors.New(cors.Config{
		AllowOrigins:     config.Cors,
		AllowMethods:     []string{"GET", "POST", "DELETE"},
		AllowHeaders:     []string{"Origin", "CONTENT-TYPE"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go serve.Run(config.App.Listen)
			return nil
		},
	})
	return
}
