package service

import (
	"context"
	jsoniter "github.com/json-iterator/go"
	"lab-api/model"
	"strconv"
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

func (x *Role) GetFromCache(ctx context.Context, roleId int64) (data map[string]interface{}, err error) {
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
	if value, err = x.Redis.HGet(ctx, x.Key, strconv.Itoa(int(roleId))).Result(); err != nil {
		return
	}
	if err = jsoniter.Unmarshal([]byte(value), &data); err != nil {
		return
	}
	return
}

func (x *Role) RefreshCache(ctx context.Context) (err error) {
	var data []struct {
		ID          int64       `json:"-"`
		Resources   model.Array `json:"resources"`
		Strategies  model.Array `json:"strategies"`
		Permissions model.Array `json:"permissions"`
	}
	if err = x.Db.WithContext(ctx).
		Table("role r").
		Select([]string{
			"r.id",
			"coalesce(json_agg(distinct rrr.resource_id) filter ( where rrr.resource_id is not null ), '[]') as resources",
			"coalesce(json_agg(distinct p.strategy) filter ( where p.strategy is not null ), '[]')           as strategies",
			"r.permissions",
		}).
		Joins("left join role_resource_rel rrr on r.id = rrr.role_id").
		Joins("left join (?) p on r.id = p.role_id", x.Db.
			Table("role_resource_rel rrr").
			Select([]string{
				"rrr.role_id",
				"array [rar.path,max(rar.mode)::char(1)] as strategy",
			}).
			Joins("join resource_acl_rel rar on rrr.resource_id = rar.resource_id").
			Group("rrr.role_id, rar.path"),
		).
		Group("r.id").
		Where("r.id <> ?", 1).
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
		values[strconv.Itoa(int(v.ID))] = value
	}
	if err = x.Redis.HMSet(ctx, x.Key, values).Err(); err != nil {
		return
	}
	return
}

func (x *Role) RemoveCache(ctx context.Context) error {
	return x.Redis.Del(ctx, x.Key).Err()
}
