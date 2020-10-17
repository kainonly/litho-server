package helper

import (
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"van-api/types"
)

var (
	Config *types.Config
	Db     *gorm.DB
	Redis  *redis.Client
)
