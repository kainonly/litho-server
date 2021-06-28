package config

type Config struct {
	Listen   string   `yaml:"listen"`
	App      App      `yaml:"app"`
	Database Database `yaml:"database"`
}

type App struct {
	Name  string `yaml:"name"`
	Key   string `yaml:"key"`
	Debug bool   `yaml:"debug"`
}

type Database struct {
	Dsn             string `yaml:"dsn"`
	MaxIdleConns    int    `yaml:"max_idle_conns"`
	MaxOpenConns    int    `yaml:"max_open_conns"`
	ConnMaxLifetime int    `yaml:"conn_max_lifetime"`
	TablePrefix     string `yaml:"table_prefix"`
}
