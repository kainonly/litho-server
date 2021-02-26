package config

import (
	"github.com/kainonly/gin-extra/cors"
	"lab-api/config/options"
)

type Config struct {
	Listen   string                 `yaml:"listen"`
	App      options.AppOption      `yaml:"app"`
	Database options.DatabaseOption `yaml:"database"`
	Redis    options.RedisOption    `yaml:"redis"`
	Cors     cors.Option            `yaml:"cors"`
	Storage  options.Storage        `yaml:"storage"`
}
