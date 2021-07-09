package service

import (
	"context"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/go-redis/redis/v8"
	jsoniter "github.com/json-iterator/go"
	"gorm.io/gorm"
	"lab-api/config"
	"lab-api/model"
)

type Acl struct {
	key   string
	db    *gorm.DB
	redis *redis.Client
}

func NewAcl(cfg *config.Config, db *gorm.DB, redis *redis.Client) *Acl {
	return &Acl{
		key:   cfg.App.Key("system:acl"),
		db:    db,
		redis: redis,
	}
}

func (x *Acl) Get(ctx context.Context, key string, mode string) (result *hashset.Set, err error) {
	var exists int64
	if exists, err = x.redis.Exists(ctx, x.key).Result(); err != nil {
		return
	}
	if exists == 0 {
		var acls []model.Acl
		if err = x.db.Where("status = ?", true).Find(&acls).Error; err != nil {
			return
		}
		lists := make(map[string]interface{})
		for _, acl := range acls {
			b, _ := jsoniter.Marshal(map[string]interface{}{
				"w": acl.Write,
				"r": acl.Read,
			})
			lists[acl.Key] = string(b)
		}
		x.redis.HMSet(ctx, x.key, lists)
	}
	result = hashset.New()
	var b []byte
	if b, err = x.redis.HGet(ctx, x.key, key).Bytes(); err != nil {
		return
	}
	var data map[string][]interface{}
	jsoniter.Unmarshal(b, &data)
	result.Add(data["r"]...)
	if mode == "1" {
		result.Add(data["w"]...)
	}
	return
}

func (x *Acl) Clear(ctx context.Context) error {
	return x.redis.Del(ctx, x.key).Err()
}
