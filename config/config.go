package config

import "taste-api/config/options"

type Config struct {
	Listen string              `yaml:"listen"`
	Debug  bool                `yaml:"debug"`
	App    options.AppOption   `yaml:"app"`
	Cors   options.CorsOption  `yaml:"cors"`
	Mysql  options.MysqlOption `yaml:"mysql"`
	Redis  options.RedisOption `yaml:"redis"`
	//Token  map[string]token.Option `yaml:"token"`
}
