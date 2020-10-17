package helper

import (
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

var (
	DB  *gorm.DB
	RDB *redis.Client
)
