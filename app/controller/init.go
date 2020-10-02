package controller

import (
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type controller struct {
	db  *gorm.DB
	rdb *redis.Client
}

func New(db *gorm.DB, rdb *redis.Client) *controller {
	c := new(controller)
	c.db = db
	c.rdb = rdb
	return c
}
