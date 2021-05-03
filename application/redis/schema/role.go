package schema

import (
	"context"
	"github.com/emirpasic/gods/sets/hashset"
	jsoniter "github.com/json-iterator/go"
	"github.com/kainonly/gin-extra/rbacx"
	"lab-api/application/model"
	"lab-api/helper"
)

type Role struct {
	rbacx.RoleAPI

	key string
	Dependency
}

func NewRole(dep Dependency) *Role {
	c := new(Role)
	c.key = "system:role"
	c.Dependency = dep
	return c
}

func (c *Role) Clear(ctx context.Context) {
	c.Redis.Del(ctx, c.key)
}

func (c *Role) Get(ctx context.Context, keys []string, mode string) *hashset.Set {
	exists := c.Redis.Exists(ctx, c.key).Val()
	if exists == 0 {
		var roleLists []model.RoleMix
		c.Db.Where("status = ?", true).Find(&roleLists)
		lists := make(map[string]interface{})
		for _, role := range roleLists {
			bs, _ := jsoniter.Marshal(map[string]interface{}{
				"acl":        helper.StringToSlice(role.Acl, ","),
				"resource":   helper.StringToSlice(role.Resource, ","),
				"permission": helper.StringToSlice(role.Permission, ","),
			})
			lists[role.Key] = string(bs)
		}
		c.Redis.HMSet(ctx, c.key, lists)
	}
	dataset := c.Redis.HMGet(ctx, c.key, keys...).Val()
	set := hashset.New()
	for _, val := range dataset {
		var data map[string]interface{}
		jsoniter.Unmarshal([]byte(val.(string)), &data)
		set.Add(data[mode].([]interface{})...)
	}
	return set
}
