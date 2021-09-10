package role

import (
	"context"
	jsoniter "github.com/json-iterator/go"
	"lab-api/model"
)

func (x *Service) GetRouters(ctx context.Context, code string) (data []int64, err error) {
	var exists int64
	if exists, err = x.Redis.Exists(ctx, x.Routers).Result(); err != nil {
		return
	}
	if exists == 0 {
		if err = x.RefreshCache(ctx); err != nil {
			return
		}
	}
	var value string
	if value, err = x.Redis.HGet(ctx, x.Routers, code).Result(); err != nil {
		return
	}
	if err = jsoniter.Unmarshal([]byte(value), &data); err != nil {
		return
	}
	return
}

func (x *Service) GetPermissions(ctx context.Context, code string) (data []string, err error) {
	var exists int64
	if exists, err = x.Redis.Exists(ctx, x.Permissions).Result(); err != nil {
		return
	}
	if exists == 0 {
		if err = x.RefreshCache(ctx); err != nil {
			return
		}
	}
	var value string
	if value, err = x.Redis.HGet(ctx, x.Permissions, code).Result(); err != nil {
		return
	}
	if err = jsoniter.Unmarshal([]byte(value), &data); err != nil {
		return
	}
	return
}

func (x *Service) RefreshCache(ctx context.Context) (err error) {
	var data []struct {
		ID          int64
		Key         string
		Routers     string
		Permissions string
	}
	if err = x.Db.WithContext(ctx).
		Model(&model.Role{}).
		Where("status = ?", true).
		Where("key <> ?", "*").
		Select([]string{"key", "routers", "permissions"}).
		Find(&data).Error; err != nil {
		return
	}
	routers := make(map[string]interface{}, len(data))
	permissions := make(map[string]interface{}, len(data))
	for _, v := range data {
		routers[v.Key] = v.Routers
		permissions[v.Key] = v.Permissions
	}
	if err = x.Redis.HMSet(ctx, x.Routers, routers).Err(); err != nil {
		return
	}
	if err = x.Redis.HMSet(ctx, x.Permissions, permissions).Err(); err != nil {
		return
	}
	return
}

func (x *Service) RemoveCache(ctx context.Context) error {
	return x.Redis.Del(ctx, x.Routers, x.Permissions).Err()
}
