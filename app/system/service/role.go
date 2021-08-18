package service

import (
	"context"
	"lab-api/model"
	"log"
)

type Role struct {
	*Dependency
	Key string
}

func NewRole(d Dependency) *Role {
	return &Role{
		Dependency: &d,
		Key:        d.Config.RedisKey("role"),
	}
}

func (x *Role) GetFromCache(ctx context.Context) (data map[string]string, err error) {
	var exists int64
	if exists, err = x.Redis.Exists(ctx, x.Key).Result(); err != nil {
		return
	}
	if exists == 0 {
		if err = x.RefreshCache(ctx); err != nil {
			return
		}
	}
	return
}

func (x *Role) RefreshCache(ctx context.Context) (err error) {
	var data []*model.Role
	if err = x.Db.Model(&model.Role{}).
		Association("Resources").
		Find(&data); err != nil {
		return
	}
	log.Println(data)
	return
}

func (x *Role) RemoveCache(ctx context.Context) error {
	return x.Redis.Del(ctx, x.Key).Err()
}
