package common

import (
	"github.com/go-redis/redis/v8"
	"github.com/kainonly/go-bit/authx"
	"github.com/kainonly/go-bit/cipher"
	"github.com/kainonly/go-bit/cookie"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type App struct {
	Name     string                 `yaml:"name"`
	Key      string                 `yaml:"key"`
	Database Database               `yaml:"database"`
	Redis    Redis                  `yaml:"redis"`
	Cookie   cookie.Option          `yaml:"cookie"`
	Cors     []string               `yaml:"cors"`
	Auth     map[string]*authx.Auth `yaml:"auth"`
}

type Database struct {
	Dsn             string `yaml:"dsn"`
	MaxIdleConns    int    `yaml:"max_idle_conns"`
	MaxOpenConns    int    `yaml:"max_open_conns"`
	ConnMaxLifetime int    `yaml:"conn_max_lifetime"`
}

type Redis struct {
	Address  string `yaml:"address"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

func (x *App) RedisKey(name string) string {
	return x.Name + ":" + name
}

type Dependency struct {
	fx.In

	App    *App
	Db     *gorm.DB
	Redis  *redis.Client
	Cookie *cookie.Cookie
	Authx  *authx.Authx
	Cipher *cipher.Cipher
}
