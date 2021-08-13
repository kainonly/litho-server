package config

import (
	"github.com/kainonly/go-bit/authx"
	"github.com/kainonly/go-bit/cookie"
)

type Config struct {
	App      App                    `yaml:"app"`
	Database Database               `yaml:"database"`
	Redis    Redis                  `yaml:"redis"`
	Cookie   cookie.Option          `yaml:"cookie"`
	Cors     map[string]CorsOption  `yaml:"cors"`
	Auth     map[string]*authx.Auth `yaml:"auth"`
}
