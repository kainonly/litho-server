package service

import (
	"context"
	jsoniter "github.com/json-iterator/go"
	"lab-api/model"
)

type Resource struct {
	*Dependent

	key string
}

func NewResource(dep *Dependent) *Resource {
	return &Resource{
		Dependent: dep,
		key:       dep.Config.App.Key("system:resource"),
	}
}

func (x *Resource) Fetch(ctx context.Context) (result []map[string]interface{}, err error) {
	var exists int64
	if exists, err = x.Redis.Exists(ctx, x.key).Result(); err != nil {
		return
	}
	if exists == 0 {
		var lists []map[string]interface{}
		if err = x.Db.Model(&model.Resource{}).
			Select([]string{"key", "parent", "name", "nav", "router", "policy", "icon"}).
			Where("status = ?", true).
			Order("sort desc").
			Scan(&lists).Error; err != nil {
			return
		}
		b, _ := jsoniter.Marshal(lists)
		if err = x.Redis.Set(ctx, x.key, string(b), 0).Err(); err != nil {
			return
		}
	}
	if b, _ := x.Redis.Get(ctx, x.key).Bytes(); b != nil {
		jsoniter.Unmarshal(b, &result)
	}
	return
}

func (x *Resource) Clear(ctx context.Context) error {
	return x.Redis.Del(ctx, x.key).Err()
}
