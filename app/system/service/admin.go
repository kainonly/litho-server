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

func NewAdmin(d Dependency) *Admin {
	return &Admin{
		Dependency: &d,
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

func (x *Admin) RefreshCache(ctx context.Context) (err error) {
	var data []struct {
		ID          uint64      `json:"id"`
		UUID        string      `json:"uuid"`
		Roles       model.Array `json:"roles"`
		Resources   model.Array `json:"resources"`
		Strategies  model.Array `json:"strategies"`
		Permissions model.Array `json:"permissions"`
	}
	if err = x.Db.WithContext(ctx).
		Table("admin a").
		Select([]string{
			"a.id", "a.uuid", "a.password",
			"json_agg(distinct arr.role_id) as roles",
			"coalesce(json_agg(distinct sarr.resource_id) filter ( where sarr.resource_id is not null ), '[]') as resources",
			"coalesce(json_agg(distinct p.strategy) filter ( where p.strategy is not null), '[]')              as strategies",
			"a.name", "a.email", "a.phone", "a.avatar", "a.permissions",
		}).
		Joins("join admin_role_rel arr on a.id = arr.admin_id").
		Joins("left join admin_resource_rel sarr on a.id = sarr.admin_id").
		Joins("left join (?) p on a.id = p.admin_id", x.Db.
			Table("admin_resource_rel arr").
			Select([]string{
				"arr.admin_id",
				"array [rar.path,max(rar.mode)::char(1)] as strategy",
			}).
			Joins("join resource_acl_rel rar on arr.resource_id = rar.resource_id").
			Group("arr.admin_id, rar.path"),
		).
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
