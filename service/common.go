package service

import (
	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
	"gorm.io/gorm"
	"lab-api/config"
)

var Provides = wire.NewSet(
	wire.Struct(new(Dependency), "*"),
	NewAdmin,
)

type Dependency struct {
	Config *config.Config
	Db     *gorm.DB
	Redis  *redis.Client
}

type Query func(tx *gorm.DB) *gorm.DB
