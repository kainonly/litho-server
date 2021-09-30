package common

import (
	"github.com/go-redis/redis/v8"
	"github.com/weplanx/support/helper"
	"github.com/weplanx/support/passport"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type App struct {
	fx.In

	Set    *Set
	Db     *gorm.DB
	Redis  *redis.Client
	Cookie *helper.CookieHelper
	Cipher *helper.CipherHelper
}

type Set struct {
	Name     string                        `yaml:"name"`
	Key      string                        `yaml:"key"`
	Database Database                      `yaml:"database"`
	Redis    Redis                         `yaml:"redis"`
	Cookie   helper.CookieOption           `yaml:"cookie"`
	Cors     []string                      `yaml:"cors"`
	Auth     map[string]*passport.Passport `yaml:"auth"`
}

type Database struct {
	Dsn             string `yaml:"dsn"`
	MaxIdleConns    int    `yaml:"max_idle_conns"`
	MaxOpenConns    int    `yaml:"max_open_conns"`
	ConnMaxLifetime int    `yaml:"conn_max_lifetime"`
	TablePrefix     string `yaml:"table_prefix"`
}

type Redis struct {
	Address  string `yaml:"address"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

func (x *Set) RedisKey(name string) string {
	return x.Name + ":" + name
}
