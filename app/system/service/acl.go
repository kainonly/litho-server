package service

import (
	"context"
	jsoniter "github.com/json-iterator/go"
	"lab-api/model"
)

type Acl struct {
	*Dependency
	Key string
}

func NewAcl(d Dependency) *Acl {
	return &Acl{
		Dependency: &d,
		Key:        d.Config.RedisKey("acl"),
	}
}

func (x *Acl) Get(ctx context.Context, key string, readonly bool) (data map[string]string, err error) {
	var exists int64
	if exists, err = x.Redis.Exists(ctx, x.Key).Result(); err != nil {
		return
	}
	if exists == 0 {
		if err = x.RefreshCache(ctx); err != nil {
			return
		}
	}
	var value string
	if value, err = x.Redis.HGet(ctx, x.Key, key).Result(); err != nil {
		return
	}
	var acts model.Acts
	if err = jsoniter.Unmarshal([]byte(value), &acts); err != nil {
		return
	}
	data = acts.R
	if !readonly {
		for k, v := range acts.W {
			data[k] = v
		}
	}
	return
}

func (x *Acl) RefreshCache(ctx context.Context) (err error) {
	var data []map[string]interface{}
	if err = x.Db.Model(&model.Acl{}).
		Select([]string{"key", "acts"}).
		Find(&data).Error; err != nil {
		return
	}
	values := make(map[string]interface{}, len(data))
	for _, v := range data {
		values[v["key"].(string)] = v["acts"]
	}
	if err = x.Redis.HMSet(ctx, x.Key, values).Err(); err != nil {
		return
	}
	return
}

func (x *Acl) RemoveCache(ctx context.Context) error {
	return x.Redis.Del(ctx, x.Key).Err()
}
