package cache

import (
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type Model struct {
	db   *gorm.DB
	rdb  *redis.Client
	keys map[string]string
}

func Initialize(db *gorm.DB, rdb *redis.Client) *Model {
	c := new(Model)
	c.db = db
	c.rdb = rdb
	c.keys = map[string]string{
		"acl":      "system:acl",
		"resource": "system:resource",
		"role":     "system:role",
		"admin":    "system:admin",
	}
	return c
}
