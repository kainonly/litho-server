package service

import (
	"github.com/go-redis/redis/v8"
	"go.uber.org/fx"
	"gorm.io/gorm"
	"lab-api/config"
)

var Provides = fx.Provide(NewAdmin)

type Dependent struct {
	fx.In

	Config *config.Config
	Db     *gorm.DB
	Redis  *redis.Client
}

type Query func(tx *gorm.DB) *gorm.DB
