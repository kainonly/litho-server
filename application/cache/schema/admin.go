package schema

import (
	"context"
	jsoniter "github.com/json-iterator/go"
	"taste-api/application/model"
)

type Admin struct {
	key string
	Dependency
}

func NewAdmin(dep Dependency) *Admin {
	c := new(Admin)
	c.key = "system:admin"
	c.Dependency = dep
	return c
}

func (c *Admin) Clear() {
	c.Redis.Del(context.Background(), c.key)
}

func (c *Admin) Get(username string) (result map[string]interface{}, err error) {
	ctx := context.Background()
	var exists int64
	exists, err = c.Redis.Exists(ctx, c.key).Result()
	if err != nil {
		return
	}
	if exists == 0 {
		var adminLists []model.Admin
		c.Db.Where("status = ?", 1).
			Find(&adminLists)

		lists := make(map[string]interface{})
		for _, admin := range adminLists {
			var bs []byte
			bs, err = jsoniter.Marshal(map[string]interface{}{
				"id":       admin.ID,
				"role":     admin.Role,
				"username": admin.Username,
				"password": admin.Password,
			})
			if err != nil {
				return
			}
			lists[admin.Username] = string(bs)
		}
		err = c.Redis.HMSet(ctx, c.key, lists).Err()
		if err != nil {
			return
		}
	}
	var bs []byte
	bs, err = c.Redis.HGet(ctx, c.key, username).Bytes()
	if err != nil {
		return
	}
	err = jsoniter.Unmarshal(bs, &result)
	if err != nil {
		return
	}
	return
}
