package bootstrap

import (
	"context"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/kainonly/gin-helper/authx"
	"github.com/kainonly/gin-helper/cookie"
	"github.com/kainonly/gin-helper/cors"
	"github.com/kainonly/gin-helper/dex"
	"go.uber.org/fx"
	"gopkg.in/yaml.v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"io/ioutil"
	"lab-api/config"
	"net/http"
	"os"
	"time"
)

var (
	LoadConfigurationNotExists = errors.New("the configuration file does not exist")
)

// LoadConfiguration 初始化应用配置
func LoadConfiguration() (cfg *config.Config, err error) {
	if _, err = os.Stat("./config.yml"); os.IsNotExist(err) {
		err = LoadConfigurationNotExists
		return
	}
	var buf []byte
	buf, err = ioutil.ReadFile("./config.yml")
	if err != nil {
		return
	}
	err = yaml.Unmarshal(buf, &cfg)
	if err != nil {
		return
	}
	return
}

// InitializeDatabase 初始化 Postgresql 数据库
// 配置文档 https://gorm.io/docs/connecting_to_the_database.html
func InitializeDatabase(cfg *config.Config) (db *gorm.DB, err error) {
	option := cfg.Database
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
func InitializeRedis(cfg *config.Config) *redis.Client {
	option := cfg.Redis
	return redis.NewClient(&redis.Options{
		Addr:     option.Address,
		Password: option.Password,
		DB:       option.DB,
	})
}

// InitializeCookie 初始化 Cookie 设置
func InitializeCookie(cfg *config.Config) *cookie.Cookie {
	return cookie.Make(cfg.Cookie, http.SameSiteStrictMode)
}

// InitializeDex 初始化数据加密
func InitializeDex(cfg *config.Config) (*dex.Dex, error) {
	return dex.Make(dex.Option{Key: cfg.App.Key})
}

// InitializeAuth 初始化鉴权
func InitializeAuth(cfg *config.Config, c *cookie.Cookie) *authx.Auth {
	option := cfg.Auth
	option.Key = cfg.App.Key
	return authx.Make(cfg.Auth, authx.Args{
		Method: jwt.SigningMethodHS256,
		UseCookie: &cookie.Cookie{
			Name:   "access_token",
			Option: c.Option,
		},
		RefreshFn: nil,
	})
}

// HttpServer 启动 Gin HTTP 服务
// 配置文档 https://gin-gonic.com/docs/examples/custom-http-config
func HttpServer(lc fx.Lifecycle, cfg *config.Config) (serve *gin.Engine) {
	serve = gin.New()
	serve.Use(gin.Logger())
	serve.Use(gin.Recovery())
	serve.Use(cors.Cors(cfg.Cors))
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go serve.Run(cfg.Listen)
			return nil
		},
	})
	return
}
