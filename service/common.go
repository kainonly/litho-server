package service

import (
	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
	"gorm.io/gorm"
)

var Provides = wire.NewSet(
	wire.Struct(new(Dependency), "*"),
	NewAdmin,
)

type Dependency struct {
	Db    *gorm.DB
	Redis *redis.Client
}

type Query func(tx *gorm.DB) *gorm.DB
