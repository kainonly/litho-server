package service

import (
	"github.com/go-redis/redis/v8"
	"go.uber.org/fx"
	"gorm.io/gorm"
	"lab-api/config"
)

type Query func(tx *gorm.DB) *gorm.DB

type Dependent struct {
	fx.In

	Config *config.Config
	Db     *gorm.DB
	Redis  *redis.Client
}

var Provides = fx.Provide(
	NewLock,
	NewAcl,
	NewResource,
	NewRole,
	NewAdmin,
)
