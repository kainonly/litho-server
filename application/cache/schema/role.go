package schema

import (
	"context"
	jsoniter "github.com/json-iterator/go"
	"strings"
	"taste-api/application/model"
)

type Role struct {
	schema
}

func NewRole(dep Dependency) *Role {
	c := new(Role)
	c.set("system:role", dep)
	return c
}

func (c *Role) RoleClear() {
	c.dep.Redis.Del(context.Background(), c.key)
}

func (c *Role) RoleGet(keys []string, mode string) (result []string, err error) {
	ctx := context.Background()
	var exists int64
	exists, err = c.dep.Redis.Exists(ctx, c.key).Result()
	if err != nil {
		return
	}
	if exists == 0 {
		var roleLists []model.Role
		c.dep.Db.Where("status = ?", 1).
			Find(&roleLists)

		lists := make(map[string]interface{})
		for _, role := range roleLists {
			var buf []byte
			buf, err = jsoniter.Marshal(map[string]interface{}{
				"acl":      role.Acl,
				"resource": role.Resource,
			})
			if err != nil {
				return
			}
			lists[role.Key] = string(buf)
		}
		err = c.dep.Redis.HMSet(ctx, c.key, lists).Err()
		if err != nil {
			return
		}
	}
	var raws []interface{}
	raws, err = c.dep.Redis.HMGet(ctx, c.key, keys...).Result()
	result = make([]string, 0)
	for _, raw := range raws {
		var value map[string]interface{}
		err = jsoniter.Unmarshal([]byte(raw.(string)), &value)
		if err != nil {
			return
		}
		result = append(result, strings.Split(value[mode].(string), ",")...)
	}
	return
}
