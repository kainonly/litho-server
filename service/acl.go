package service

import (
	"context"
	"github.com/emirpasic/gods/sets/hashset"
	jsoniter "github.com/json-iterator/go"
	"github.com/kainonly/gin-helper/rbac"
	"lab-api/model"
)

type Acl struct {
	rbac.AclFn
	*Dependent

	key string
}

func NewAcl(dep *Dependent) *Acl {
	return &Acl{
		key:       dep.Config.App.Key("system:acl"),
		Dependent: dep,
	}
}

func (x *Acl) Fetch(ctx context.Context, key string, mode string) (result *hashset.Set, err error) {
	var exists int64
	if exists, err = x.Redis.Exists(ctx, x.key).Result(); err != nil {
		return
	}
	if exists == 0 {
		var acls []model.Acl
		if err = x.Db.Where("status = ?", true).Find(&acls).Error; err != nil {
			return
		}
		lists := make(map[string]interface{})
		for _, acl := range acls {
			b, _ := jsoniter.Marshal(map[string]interface{}{
				"w": acl.Write,
				"r": acl.Read,
			})
			lists[acl.Key] = string(b)
		}
		x.Redis.HMSet(ctx, x.key, lists)
	}
	result = hashset.New()
	var b []byte
	if b, err = x.Redis.HGet(ctx, x.key, key).Bytes(); err != nil {
		return
	}
	var data map[string][]interface{}
	jsoniter.Unmarshal(b, &data)
	result.Add(data["r"]...)
	if mode == "1" {
		result.Add(data["w"]...)
	}
	return
}

func (x *Acl) Clear(ctx context.Context) error {
	return x.Redis.Del(ctx, x.key).Err()
}
