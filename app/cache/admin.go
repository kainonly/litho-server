package cache

import (
	"context"
	jsoniter "github.com/json-iterator/go"
	"van-api/app/model"
)

func (c *Model) AdminClear() {
	c.rdb.Del(context.Background(), c.keys["admin"])
}

func (c *Model) AdminGet(username string) (result map[string]interface{}, err error) {
	ctx := context.Background()
	var exists int64
	exists, err = c.rdb.Exists(ctx, c.keys["admin"]).Result()
	if err != nil {
		return
	}
	if exists == 0 {
		var adminLists []model.Admin
		c.db.Where("status = ?", 1).
			Find(&adminLists)

		lists := make(map[string]interface{})
		for _, data := range adminLists {
			var buf []byte
			value := map[string]interface{}{
				"id":       data.ID,
				"role":     data.Role,
				"username": data.Username,
				"password": data.Password,
			}
			buf, err = jsoniter.Marshal(value)
			if err != nil {
				return
			}
			if data.Username == username {
				result = value
			}
			lists[data.Username] = string(buf)
		}
		err = c.rdb.HMSet(ctx, c.keys["admin"], lists).Err()
		if err != nil {
			return
		}
	} else {
		var raw []byte
		err = c.rdb.HGet(ctx, c.keys["admin"], username).Scan(&raw)
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
