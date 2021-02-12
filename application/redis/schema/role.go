package schema

import (
	"context"
	"github.com/emirpasic/gods/sets/hashset"
	jsoniter "github.com/json-iterator/go"
	"lab-api/application/model"
)

type Role struct {
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
				"acl":        role.Acl,
				"resource":   role.Resource,
				"permission": role.Permission,
			})
			lists[role.Key] = string(bs)
		}
		c.Redis.HMSet(ctx, c.key, lists)
	}
	lists := c.Redis.HMGet(ctx, c.key, keys...).Val()
	set := hashset.New()
	for _, val := range lists {
		if val != nil {
			var data map[string]interface{}
			jsoniter.Unmarshal([]byte(val.(string)), &data)
			set.Add(data[mode].([]interface{})...)
		}
	}
	return set
}
