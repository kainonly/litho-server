package schema

import (
	"context"
	"github.com/emirpasic/gods/sets/hashset"
	jsoniter "github.com/json-iterator/go"
	"github.com/kainonly/gin-extra/rbacx"
	"lab-api/application/model"
	"strings"
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

func (c *Role) Clear() {
	c.Redis.Del(context.Background(), c.key)
}

func (c *Role) Get(keys []string, mode string) *hashset.Set {
	ctx := context.Background()
	exists := c.Redis.Exists(ctx, c.key).Val()
	if exists == 0 {
		var roleLists []model.RoleMix
		c.Db.Where("status = ?", true).Find(&roleLists)
		lists := make(map[string]interface{})
		for _, role := range roleLists {
			bs, _ := jsoniter.Marshal(map[string]interface{}{
				"acl":        strings.Split(role.Acl, ","),
				"resource":   strings.Split(role.Resource, ","),
				"permission": strings.Split(role.Permission, ","),
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
