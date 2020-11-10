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
