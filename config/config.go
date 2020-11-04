package config

import (
	"taste-api/config/options"
)

type Config struct {
	Listen   string                         `yaml:"listen"`
	App      options.AppOption              `yaml:"app"`
	Database options.DatabaseOption         `yaml:"database"`
	Redis    options.RedisOption            `yaml:"redis"`
	Cors     options.CorsOption             `yaml:"cors"`
	Token    map[string]options.TokenOption `yaml:"token"`
}
