package common

import (
	"github.com/kainonly/go-bit/authx"
	"github.com/kainonly/go-bit/cookie"
)

type Config struct {
	App      App
	Database Database
	Redis    Redis
	Cookie   cookie.Option
	Cors     []string `env:"CORS" envSeparator:","`
	Auth     map[string]*authx.Auth
}

type App struct {
	Name string `env:"APP_NAME"`
	Key  string `env:"APP_KEY"`
}

type Database struct {
	Dsn             string `env:"DB_DSN"`
	MaxIdleConns    int    `env:"DB_MAX_IDLE_CONNS" envDefault:"5"`
	MaxOpenConns    int    `env:"DB_MAX_OPEN_CONNS" envDefault:"10"`
	ConnMaxLifetime int    `env:"DB_CONN_MAX_LIFETIME" envDefault:"3600"`
}

type Redis struct {
	Address  string `env:"REDIS_ADDRESS"`
	Password string `env:"REDIS_PASSWORD"`
	DB       int    `env:"REDIS_DB"`
}

func (x *Config) RedisKey(name string) string {
	return x.App.Name + ":" + name
}
