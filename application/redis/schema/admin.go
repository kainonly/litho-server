package schema

import (
	"context"
	jsoniter "github.com/json-iterator/go"
	"github.com/kainonly/gin-extra/rbacx"
	"lab-api/application/model"
	"lab-api/helper"
)

type Admin struct {
	rbacx.UserAPI

	key string
	Dependency
}

func NewAdmin(dep Dependency) *Admin {
	c := new(Admin)
	c.key = "system:admin"
	c.Dependency = dep
	return c
}

func (c *Admin) Clear(ctx context.Context) {
	c.Redis.Del(ctx, c.key)
}

func (c *Admin) Get(ctx context.Context, username string) (result map[string]interface{}) {
	exists := c.Redis.Exists(ctx, c.key).Val()
	if exists == 0 {
		var adminLists []model.AdminMix
		c.Db.Where("status = ?", true).Find(&adminLists)

		lists := make(map[string]interface{})
		for _, admin := range adminLists {
			bs, _ := jsoniter.Marshal(map[string]interface{}{
				"id":         admin.ID,
				"role":       helper.StringToSlice(admin.Role, ","),
				"username":   admin.Username,
				"password":   admin.Password,
				"resource":   helper.StringToSlice(admin.Resource, ","),
				"acl":        helper.StringToSlice(admin.Acl, ","),
				"permission": helper.StringToSlice(admin.Permission, ","),
			})
			lists[admin.Username] = string(bs)
		}
		c.Redis.HMSet(ctx, c.key, lists)
	}
	if bs, _ := c.Redis.HGet(ctx, c.key, username).Bytes(); bs != nil {
		jsoniter.Unmarshal(bs, &result)
	}
	return
}
