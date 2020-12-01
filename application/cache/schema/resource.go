package schema

import (
	"context"
	jsoniter "github.com/json-iterator/go"
	"lab-api/application/model"
)

type Resource struct {
	key string
	Dependency
}

func NewResource(dep Dependency) *Resource {
	c := new(Resource)
	c.key = "system:resource"
	c.Dependency = dep
	return c
}

func (c *Resource) Clear() {
	c.Redis.Del(context.Background(), c.key)
}

func (c *Resource) Get() (result []map[string]interface{}) {
	ctx := context.Background()
	exists := c.Redis.Exists(ctx, c.key).Val()
	if exists == 0 {
		var resourceLists []map[string]interface{}
		c.Db.Model(&model.Resource{}).
			Select([]string{"key", "parent", "name", "nav", "router", "policy", "icon"}).
			Where("status = ?", true).
			Order("sort desc").
			Scan(&resourceLists)
		if len(resourceLists) != 0 {
			bs, _ := jsoniter.Marshal(resourceLists)
			c.Redis.Set(ctx, c.key, string(bs), 0)
		}

	}
	if bs, _ := c.Redis.Get(ctx, c.key).Bytes(); bs != nil {
		jsoniter.Unmarshal(bs, &result)
	}
	return
}
