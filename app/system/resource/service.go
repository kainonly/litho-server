package resource

import (
	"context"
	jsoniter "github.com/json-iterator/go"
	"lab-api/model"
)

func (x *Service) GetFromCache(ctx context.Context) (data []model.Resource, err error) {
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

func (x *Service) RefreshCache(ctx context.Context) (err error) {
	var data []model.Resource
	if err = x.Db.WithContext(ctx).
		Model(&model.Resource{}).
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

func (x *Service) RemoveCache(ctx context.Context) error {
	return x.Redis.Del(ctx, x.Key).Err()
}
