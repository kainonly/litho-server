package service

import (
	"context"
	jsoniter "github.com/json-iterator/go"
	"lab-api/model"
)

type Role struct {
	*Dependency
	Key string
}

func NewRole(d *Dependency) *Role {
	return &Role{
		Dependency: d,
		Key:        d.Config.RedisKey("role"),
	}
}

func (x *Role) GetFromCache(ctx context.Context, code string) (data map[string]interface{}, err error) {
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
	if value, err = x.Redis.HGet(ctx, x.Key, code).Result(); err != nil {
		return
	}
	if err = jsoniter.Unmarshal([]byte(value), &data); err != nil {
		return
	}
	return
}

func (x *Role) RefreshCache(ctx context.Context) (err error) {
	var data []struct {
		ID        int64       `json:"-"`
		Code      string      `json:"code"`
		Resources model.Array `json:"resources"`
	}
	if err = x.Db.WithContext(ctx).
		Table("role r").
		Select([]string{
			"r.id", "r.code",
			"json_agg(rrr.resource_path) as resources",
		}).
		Joins("join role_resource_rel rrr on r.code = rrr.role_code").
		Group("r.id").
		Where("r.code <> ?", "*").
		Where("r.status = ?", true).
		Scan(&data).Error; err != nil {
		return
	}
	values := make(map[string]interface{}, len(data))
	for _, v := range data {
		var value []byte
		if value, err = jsoniter.Marshal(&v); err != nil {
			return
		}
		values[v.Code] = value
	}
	if err = x.Redis.HMSet(ctx, x.Key, values).Err(); err != nil {
		return
	}
	return
}

func (x *Role) RemoveCache(ctx context.Context) error {
	return x.Redis.Del(ctx, x.Key).Err()
}
