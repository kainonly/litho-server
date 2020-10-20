package cache

import (
	"context"
	jsoniter "github.com/json-iterator/go"
)

type Resource struct {
	Key    string `json:"key"`
	Parent string `json:"parent"`
	Name   string `json:"name"`
	Nav    bool   `json:"nav"`
	Router bool   `json:"router"`
	Policy uint8  `json:"policy"`
	Icon   string `json:"icon"`
}

func (c *Model) ResourceClear() {
	c.rdb.Del(context.Background(), c.keys["resource"])
}

func (c *Model) ResourceGet() (result []Resource, err error) {
	ctx := context.Background()
	var exists int64
	exists, err = c.rdb.Exists(ctx, c.keys["resource"]).Result()
	if err != nil {
		return
	}
	if exists == 0 {
		c.db.Table("resource").
			Select([]string{"`key`", "parent", "name", "nav", "router", "policy", "icon"}).
			Where("status = ?", 1).
			Order("sort").
			Scan(&result)
		var buf []byte
		buf, err = jsoniter.Marshal(result)
		if err != nil {
			return
		}
		err = c.rdb.Set(ctx, c.keys["resource"], string(buf), 0).Err()
		if err != nil {
			return
		}
	} else {
		var raw []byte
		err = c.rdb.Get(ctx, c.keys["resource"]).Scan(&raw)
		if err != nil {
			return
		}
		err = jsoniter.Unmarshal(raw, &result)
		if err != nil {
			return
		}
	}
	return
}
