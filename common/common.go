package common

import (
	"github.com/kainonly/go-bit/authx"
	"github.com/kainonly/go-bit/cookie"
)

type Config struct {
	App      App                    `yaml:"app"`
	Database Database               `yaml:"database"`
	Redis    Redis                  `yaml:"redis"`
	Cookie   cookie.Option          `yaml:"cookie"`
	Cors     []string               `yaml:"cors"`
	Auth     map[string]*authx.Auth `yaml:"auth"`
}

type App struct {
	Name string `yaml:"name"`
	Key  string `yaml:"key"`
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

func (x *Config) RedisKey(name string) string {
	return x.App.Name + ":" + name
}
