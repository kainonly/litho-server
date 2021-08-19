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
		ID          int64           `json:"-"`
		Resources   model.JSONArray `json:"resources"`
		Acls        model.JSONArray `json:"acls"`
		Permissions model.JSONArray `json:"permissions"`
	}
	if err = x.Db.WithContext(ctx).
		Table("role r").
		Select([]string{
			"r.id",
			"json_agg(distinct rrr.resource_id) as resources",
			"json_agg(distinct ap.policy) as acls",
			"json_agg(distinct pp.code) as permissions",
		}).
		Joins("left join role_resource_rel rrr on r.id = rrr.role_id").
		Joins("left join (?) ap on r.id = ap.role_id", x.Db.
			Table("role_resource_rel rrr").
			Select([]string{
				"rrr.role_id",
				"json_build_array(a.model, max(rar.policy))::jsonb as policy",
			}).
			Joins("join resource_acl_rel rar on rrr.resource_id = rar.resource_id").
			Joins("left join acl a on rar.acl_id = a.id").
			Group("rrr.role_id, a.model"),
		).
		Joins("left join (?) pp on r.id = pp.role_id", x.Db.
			Table("role_permission_rel rpr").
			Joins("left join permission p on rpr.permission_id = p.id"),
		).
		Group("r.id,r.name").
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
