package service

import (
	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
	"gorm.io/gorm"
	"lab-api/common"
)

type Dependency struct {
	Config common.Config
	Db     *gorm.DB
	Redis  *redis.Client
}

var Provides = wire.NewSet(
	wire.Struct(new(Dependency), "*"),
	NewIndex,
)
