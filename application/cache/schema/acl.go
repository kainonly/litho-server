package schema

import (
	"context"
	jsoniter "github.com/json-iterator/go"
	"strings"
	"taste-api/application/model"
)

type Acl struct {
	key string
	dep Dependency
}

func NewAcl(dep Dependency) *Acl {
	c := new(Acl)
	c.key = "system:acl"
	c.dep = dep
	return c
}

func (c *Acl) Clear() {
	c.dep.Redis.Del(context.Background(), c.key)
}

func (c *Acl) Get(key string, policy uint8) (result []string, err error) {
	ctx := context.Background()
	var exists int64
	exists, err = c.dep.Redis.Exists(ctx, c.key).Result()
	if err != nil {
		return
	}
	if exists == 0 {
		var aclLists []model.Acl
		c.dep.Db.Where("status = ?", 1).
			Find(&aclLists)

		lists := make(map[string]interface{})
		for _, acl := range aclLists {
			var buf []byte
			buf, err = jsoniter.Marshal(map[string]interface{}{
				"write": acl.Write,
				"read":  acl.Read,
			})
			if err != nil {
				return
			}
			lists[acl.Key] = string(buf)
		}
		err = c.dep.Redis.HMSet(ctx, c.key, lists).Err()
		if err != nil {
			return
		}
	}
	var raw []byte
	raw, err = c.dep.Redis.HGet(ctx, c.key, key).Bytes()
	if err != nil {
		return
	}
	var data map[string]interface{}
	err = jsoniter.Unmarshal(raw, &data)
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
