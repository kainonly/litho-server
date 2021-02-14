package schema

import (
	"context"
	jsoniter "github.com/json-iterator/go"
	"github.com/kainonly/gin-extra/rbacx"
	"lab-api/application/model"
	"strings"
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

func (c *Admin) Clear() {
	c.Redis.Del(context.Background(), c.key)
}

func (c *Admin) Get(username string) (result map[string]interface{}) {
	ctx := context.Background()
	exists := c.Redis.Exists(ctx, c.key).Val()
	if exists == 0 {
		var adminLists []model.AdminMix
		c.Db.Where("status = ?", true).Find(&adminLists)

		lists := make(map[string]interface{})
		for _, admin := range adminLists {
			bs, _ := jsoniter.Marshal(map[string]interface{}{
				"id":         admin.ID,
				"role":       admin.Role,
				"username":   admin.Username,
				"password":   admin.Password,
				"resource":   strings.Split(admin.Resource, ","),
				"acl":        strings.Split(admin.Acl, ","),
				"permission": strings.Split(admin.Permission, ","),
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
