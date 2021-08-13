package config

type App struct {
	Listen string `yaml:"listen"`
	Name   string `yaml:"name"`
	Key    string `yaml:"key"`
}
