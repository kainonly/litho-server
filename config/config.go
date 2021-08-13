package config

import (
	"github.com/kainonly/go-bit/authx"
	"github.com/kainonly/go-bit/cookie"
)

type Config struct {
	App      App           `yaml:"app"`
	Database Database      `yaml:"database"`
	Redis    Redis         `yaml:"redis"`
	Cors     CorsOption    `yaml:"cors"`
	Cookie   cookie.Option `yaml:"cookie"`
	Auth     authx.Option  `yaml:"auth"`
}
