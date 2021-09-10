package admin

import (
	"context"
	jsoniter "github.com/json-iterator/go"
	"lab-api/model"
)

func (x *Service) FindByUsername(ctx context.Context, username string) (data model.Admin, err error) {
	if err = x.Db.WithContext(ctx).
		Where("username = ?", username).
		Where("status = ?", true).
		First(&data).Error; err != nil {
		return
	}
	return
}

func (x *Service) GetFromCache(ctx context.Context, uuid string) (data map[string]interface{}, err error) {
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
	if value, err = x.Redis.HGet(ctx, x.Key, uuid).Result(); err != nil {
		return
	}
	if err = jsoniter.Unmarshal([]byte(value), &data); err != nil {
		return
	}
	return
}

func (x *Service) RefreshCache(ctx context.Context) (err error) {
	var data []model.Admin
	if err = x.Db.WithContext(ctx).
		Find(&data).Error; err != nil {
		return
	}
	values := make(map[string]interface{}, len(data))
	for _, v := range data {
		var value []byte
		if value, err = jsoniter.Marshal(&map[string]interface{}{
			"roles":       v.Roles,
			"routers":     v.Routers,
			"permissions": v.Permissions,
		}); err != nil {
			return
		}
		values[v.Uuid.String()] = value
	}
	if err = x.Redis.HMSet(ctx, x.Key, values).Err(); err != nil {
		return
	}
	return
}

func (x *Service) RemoveCache(ctx context.Context) error {
	return x.Redis.Del(ctx, x.Key).Err()
}
