package schema

import (
	"context"
	"github.com/emirpasic/gods/sets/hashset"
	jsoniter "github.com/json-iterator/go"
	"github.com/kainonly/gin-extra/rbacx"
	"github.com/kainonly/gin-extra/str"
	"lab-api/application/model"
	"strings"
)

type Acl struct {
	rbacx.AclAPI

	key string
	Dependency
}

func NewAcl(dep Dependency) *Acl {
	c := new(Acl)
	c.key = "system:acl"
	c.Dependency = dep
	return c
}

func (c *Acl) Clear() {
	c.Redis.Del(context.Background(), c.key)
}

func (c *Acl) Get(key string, policy string) *hashset.Set {
	ctx := context.Background()
	exists := c.Redis.Exists(ctx, c.key).Val()
	if exists == 0 {
		var aclLists []model.Acl
		c.Db.Where("status = ?", true).Find(&aclLists)

		lists := make(map[string]interface{})
		for _, acl := range aclLists {
			bs, _ := jsoniter.Marshal(map[string]interface{}{
				"write": str.FilterEmpty(strings.Split(acl.Write, ",")),
				"read":  str.FilterEmpty(strings.Split(acl.Read, ",")),
			})
			lists[acl.Key] = string(bs)
		}
		c.Redis.HMSet(ctx, c.key, lists)
	}
	set := hashset.New()
	if bs, _ := c.Redis.HGet(ctx, c.key, key).Bytes(); bs != nil {
		var data map[string][]interface{}
		jsoniter.Unmarshal(bs, &data)
		set.Add(data["read"]...)
		if policy == "1" {
			set.Add(data["write"]...)
		}
	}
	return set
}
