package cache

import (
	"context"
	jsoniter "github.com/json-iterator/go"
	"strings"
	"taste-api/application/model"
)

func (c *Model) RoleClear() {
	c.rdb.Del(context.Background(), c.keys["role"])
}

func (c *Model) RoleGet(keys []string, mode string) (result []string, err error) {
	ctx := context.Background()
	var exists int64
	exists, err = c.rdb.Exists(ctx, c.keys["role"]).Result()
	if err != nil {
		return
	}
	if exists == 0 {
		var roleLists []model.Role
		c.db.Where("status = ?", 1).
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
		err = c.rdb.HMSet(ctx, c.keys["role"], lists).Err()
		if err != nil {
			return
		}
	}
	var raws []interface{}
	raws, err = c.rdb.HMGet(ctx, c.keys["role"], keys...).Result()
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
