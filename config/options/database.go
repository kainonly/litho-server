package options

import (
	"gorm.io/gorm/clause"
)

type DatabaseOption struct {
	Dsn             string `yaml:"dsn"`
	MaxIdleConns    int    `yaml:"max_idle_conns"`
	MaxOpenConns    int    `yaml:"max_open_conns"`
	ConnMaxLifetime int    `yaml:"conn_max_lifetime"`
	TablePrefix     string `yaml:"table_prefix"`
}

func (c *DatabaseOption) Table(name string) interface{} {
	return clause.Table{Name: c.TablePrefix + name, Raw: true}
}
