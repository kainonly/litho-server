package cache

import (
	"context"
	jsoniter "github.com/json-iterator/go"
	"strings"
	"taste-api/application/model"
)

func (c *Model) AclClear() {
	c.rdb.Del(context.Background(), c.keys["acl"])
}

func (c *Model) AclGet(key string, policy uint8) (result []string, err error) {
	ctx := context.Background()
	var exists int64
	exists, err = c.rdb.Exists(ctx, c.keys["acl"]).Result()
	if err != nil {
		return
	}
	if exists == 0 {
		var aclLists []model.Acl
		c.db.Where("status = ?", 1).
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
		err = c.rdb.HMSet(ctx, c.keys["acl"], lists).Err()
		if err != nil {
			return
		}
	}
	var raw []byte
	raw, err = c.rdb.HGet(ctx, c.keys["acl"], key).Bytes()
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
