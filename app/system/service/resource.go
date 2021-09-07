package service

import (
	"context"
	jsoniter "github.com/json-iterator/go"
	"github.com/kainonly/go-bit/support"
)

type Resource struct {
	*Dependency
	Key string
}

func NewResource(d Dependency) *Resource {
	return &Resource{
		Dependency: &d,
		Key:        d.Config.RedisKey("resource"),
	}
}
func (x *Resource) GetFromCache(ctx context.Context) (data []map[string]interface{}, err error) {
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
	if value, err = x.Redis.Get(ctx, x.Key).Result(); err != nil {
		return
	}
	if err = jsoniter.Unmarshal([]byte(value), &data); err != nil {
		return
	}
	return
}

func (x *Resource) RefreshCache(ctx context.Context) (err error) {
	var data []map[string]interface{}
	if err = x.Db.WithContext(ctx).
		Model(&support.Resource{}).
		Omit("status,create_time,update_time").
		Where("status = ?", true).
		Find(&data).Error; err != nil {
		return
	}
	var value []byte
	if value, err = jsoniter.Marshal(&data); err != nil {
		return
	}
	if err = x.Redis.Set(ctx, x.Key, value, 0).Err(); err != nil {
		return
	}
	return
}

func (x *Resource) RemoveCache(ctx context.Context) error {
	return x.Redis.Del(ctx, x.Key).Err()
}
