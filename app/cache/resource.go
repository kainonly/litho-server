package cache

import (
	"context"
	jsoniter "github.com/json-iterator/go"
	"van-api/app/model"
)

func (c *Model) ResourceClear() {
	c.rdb.Del(context.Background(), c.keys["resource"])
}

func (c *Model) ResourceGet() (result []map[string]interface{}, err error) {
	ctx := context.Background()
	var exists int64
	exists, err = c.rdb.Exists(ctx, c.keys["resource"]).Result()
	if err != nil {
		return
	}
	if exists == 0 {
		var resourceLists []map[string]interface{}
		c.db.Model(&model.Resource{}).
			Select([]string{"keyid", "parent", "name", "nav", "router", "policy", "icon"}).
			Where("status = ?", 1).
			Order("sort").
			Scan(&resourceLists)
		var buf []byte
		buf, err = jsoniter.Marshal(resourceLists)
		if err != nil {
			return
		}
		err = c.rdb.Set(ctx, c.keys["resource"], string(buf), 0).Err()
		if err != nil {
			return
		}
	}
	var raw []byte
	raw, err = c.rdb.Get(ctx, c.keys["resource"]).Bytes()
	if err != nil {
		return
	}
	err = jsoniter.Unmarshal(raw, &result)
	if err != nil {
		return
	}
	return
}
