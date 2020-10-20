package redis_model

import (
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type RedisModel struct {
	db  *gorm.DB
	rdb *redis.Client
}

func NewRedisModel(db *gorm.DB, rdb *redis.Client) *RedisModel {
	c := new(RedisModel)
	c.db = db
	c.rdb = rdb
	return c
}
