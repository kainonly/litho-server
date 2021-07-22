package config

import (
	"github.com/kainonly/gin-helper/authx"
	"github.com/kainonly/gin-helper/cookie"
	"github.com/kainonly/gin-helper/cors"
)

type Config struct {
	Listen   string        `yaml:"listen"`
	App      App           `yaml:"app"`
	Database Database      `yaml:"database"`
	Redis    Redis         `yaml:"redis"`
	Auth     authx.Option  `yaml:"auth"`
	Cookie   cookie.Option `yaml:"cookie"`
	Cors     cors.Option   `yaml:"cors"`
}

type App struct {
	Debug bool   `yaml:"debug"`
	Name  string `yaml:"name"`
	Key   string `yaml:"key"`
	Lock  Lock   `yaml:"lock"`
}

type Lock struct {
	Limit        int64 `yaml:"limit"`
	RecoveryTime int64 `yaml:"recovery_time"`
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
