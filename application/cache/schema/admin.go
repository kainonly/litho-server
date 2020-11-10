package schema

import (
	"context"
	jsoniter "github.com/json-iterator/go"
	"taste-api/application/model"
)

type Admin struct {
	key string
	dep Dependency
}

func NewAdmin(dep Dependency) *Admin {
	c := new(Admin)
	c.key = "system:admin"
	c.dep = dep
	return c
}

func (c *Admin) Clear() {
	c.dep.Redis.Del(context.Background(), c.key)
}

func (c *Admin) Get(username string) (result map[string]interface{}, err error) {
	ctx := context.Background()
	var exists int64
	exists, err = c.dep.Redis.Exists(ctx, c.key).Result()
	if err != nil {
		return
	}
	if exists == 0 {
		var adminLists []model.Admin
		c.dep.Db.Where("status = ?", 1).
			Find(&adminLists)

		lists := make(map[string]interface{})
		for _, admin := range adminLists {
			var buf []byte
			buf, err = jsoniter.Marshal(map[string]interface{}{
				"id":       admin.ID,
				"role":     admin.Role,
				"username": admin.Username,
				"password": admin.Password,
			})
			if err != nil {
				return
			}
			lists[admin.Username] = string(buf)
		}
		err = c.dep.Redis.HMSet(ctx, c.key, lists).Err()
		if err != nil {
			return
		}
	}
	var raw []byte
	raw, err = c.dep.Redis.HGet(ctx, c.key, username).Bytes()
	if err != nil {
		return
	}
	err = jsoniter.Unmarshal(raw, &result)
	if err != nil {
		return
	}
	return
}
