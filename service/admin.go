package service

import (
	"context"
	jsoniter "github.com/json-iterator/go"
	"github.com/kainonly/gin-helper/hash"
	"github.com/kainonly/gin-helper/rbac"
	"lab-api/model"
	"strconv"
)

type Admin struct {
	rbac.UserFn
	*Dependent

	key string
}

func NewAdmin(dep *Dependent) *Admin {
	return &Admin{
		Dependent: dep,
		key:       dep.Config.App.Key("system:admin"),
	}
}

func (x *Admin) FindOne(query Query) (data model.Admin, err error) {
	if err = query(x.Db).First(&data).Error; err != nil {
		return
	}
	return
}

func (x *Admin) Data(admin model.Admin) model.Admin {
	var password string
	if admin.Password != "" {
		password, _ = hash.Make(admin.Password)
	}
	return model.Admin{
		Email:    admin.Email,
		Password: password,
		Name:     admin.Name,
		Status:   admin.Status,
	}
}

func (x *Admin) Fetch(ctx context.Context, uid interface{}) (result map[string]interface{}, err error) {
	var exists int64
	if exists, err = x.Redis.Exists(ctx, x.key).Result(); err != nil {
		return
	}
	if exists == 0 {
		var admins []model.AdminMix
		if err = x.Db.Where("status = ?", true).Find(&admins).Error; err != nil {
			return
		}
		lists := make(map[string]interface{})
		for _, admin := range admins {
			b, _ := jsoniter.Marshal(map[string]interface{}{
				"id":         admin.ID,
				"username":   admin.Username,
				"password":   admin.Password,
				"role":       admin.Role,
				"resource":   admin.Resource,
				"acl":        admin.Acl,
				"permission": admin.Permission,
			})
			lists[strconv.FormatUint(admin.ID, 10)] = string(b)
		}
		x.Redis.HMSet(ctx, x.key, lists)
	}
	var b []byte
	if b, err = x.Redis.HGet(ctx, x.key, uid.(string)).Bytes(); err != nil {
		return
	}
	if b != nil {
		jsoniter.Unmarshal(b, &result)
	}
	return
}

func (x *Admin) Clear(ctx context.Context) error {
	return x.Redis.Del(ctx, x.key).Err()
}
