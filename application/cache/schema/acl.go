package schema

import (
	"context"
	jsoniter "github.com/json-iterator/go"
	"strings"
	"taste-api/application/model"
)

type Acl struct {
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

func (c *Acl) Get(key string, policy uint8) (result []string, err error) {
	ctx := context.Background()
	var exists int64
	exists, err = c.Redis.Exists(ctx, c.key).Result()
	if err != nil {
		return
	}
	if exists == 0 {
		var aclLists []model.Acl
		c.Db.Where("status = ?", 1).
			Find(&aclLists)

		lists := make(map[string]interface{})
		for _, acl := range aclLists {
			var bs []byte
			bs, err = jsoniter.Marshal(map[string]interface{}{
				"write": acl.Write,
				"read":  acl.Read,
			})
			if err != nil {
				return
			}
			lists[acl.Key] = string(bs)
		}
		err = c.Redis.HMSet(ctx, c.key, lists).Err()
		if err != nil {
			return
		}
	}
	var bs []byte
	bs, err = c.Redis.HGet(ctx, c.key, key).Bytes()
	if err != nil {
		return
	}
	var data map[string]interface{}
	err = jsoniter.Unmarshal(bs, &data)
	if err != nil {
		return
	}
	if policy == 0 {
		result = strings.Split(data["read"].(string), ",")
	}
	if policy == 1 {
		result = append(
			strings.Split(data["read"].(string), ","),
			strings.Split(data["write"].(string), ",")...,
		)
	}
	return
}
