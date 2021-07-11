package service

import (
	"context"
	"github.com/emirpasic/gods/sets/hashset"
	jsoniter "github.com/json-iterator/go"
	"github.com/kainonly/gin-helper/rbac"
	"lab-api/model"
)

type Role struct {
	rbac.RoleFn
	*Dependent

	key string
}

func NewRole(dep *Dependent) *Role {
	return &Role{
		Dependent: dep,
		key:       dep.Config.App.Key("system:role"),
	}
}

func (x *Role) Fetch(ctx context.Context, keys []string, mode rbac.RoleMode) (result *hashset.Set, err error) {
	var exists int64
	if exists, err = x.Redis.Exists(ctx, x.key).Result(); err != nil {
		return
	}
	if exists == 0 {
		var roles []model.RoleMix
		if err = x.Db.Where("status = ?", true).Find(&roles).Error; err != nil {
			return
		}
		lists := make(map[string]interface{})
		for _, role := range roles {
			b, _ := jsoniter.Marshal(map[string]interface{}{
				string(rbac.RoleAcl):        role.Acl,
				string(rbac.RoleResource):   role.Resource,
				string(rbac.RolePermission): role.Permission,
			})
			lists[role.Key] = string(b)
		}
		x.Redis.HMSet(ctx, x.key, lists)
	}
	result = hashset.New()
	var data []interface{}
	if data, err = x.Redis.HMGet(ctx, x.key, keys...).Result(); err != nil {
		return
	}
	for _, val := range data {
		var data map[string]interface{}
		jsoniter.Unmarshal([]byte(val.(string)), &data)
		result.Add(data[string(mode)].([]interface{})...)
	}
	return
}

func (x *Role) Clear(ctx context.Context) error {
	return x.Redis.Del(ctx, x.key).Err()
}
