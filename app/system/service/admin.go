package service

import (
	"context"
	jsoniter "github.com/json-iterator/go"
	"lab-api/model"
)

type Admin struct {
	*Dependency
	Key string
}

func NewAdmin(d *Dependency) *Admin {
	return &Admin{
		Dependency: d,
		Key:        d.Config.RedisKey("admin"),
	}
}

func (x *Admin) FindByUsername(ctx context.Context, username string) (data model.Admin, err error) {
	if err = x.Db.WithContext(ctx).
		Where("username = ?", username).
		Where("status = ?", true).
		First(&data).Error; err != nil {
		return
	}
	return
}

func (x *Admin) GetFromCache(ctx context.Context, uuid string) (data map[string]interface{}, err error) {
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

func (x *Admin) RefreshCache(ctx context.Context) (err error) {
	var data []struct {
		ID        uint64      `json:"-"`
		UUID      string      `json:"uuid"`
		Roles     model.Array `json:"roles"`
		Resources model.Array `json:"resources"`
	}
	if err = x.Db.WithContext(ctx).
		Table("admin a").
		Select([]string{
			"a.id", "a.uuid",
			"json_agg(distinct r.role_code) as roles",
			"coalesce(json_agg(distinct rr.resource_path) filter ( where rr.resource_path is not null ), '[]') as resources",
		}).
		Joins("join admin_role_rel r on a.id = r.admin_id").
		Joins("left join admin_resource_rel rr on a.id = rr.admin_id").
		Group("a.id").
		Where("status = ?", true).
		Find(&data).Error; err != nil {
		return
	}
	values := make(map[string]interface{}, len(data))
	for _, v := range data {
		var value []byte
		if value, err = jsoniter.Marshal(&v); err != nil {
			return
		}
		values[v.UUID] = value
	}
	if err = x.Redis.HMSet(ctx, x.Key, values).Err(); err != nil {
		return
	}
	return
}

func (x *Admin) RemoveCache(ctx context.Context) error {
	return x.Redis.Del(ctx, x.Key).Err()
}
