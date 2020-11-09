package schema

import (
	"github.com/go-redis/redis/v8"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type Dependency struct {
	fx.In

	Db    *gorm.DB
	Redis *redis.Client
}

type schema struct {
	key string
	dep Dependency
}

func (c *schema) set(key string, dep Dependency) {
	c.key = key
	c.dep = dep
}
