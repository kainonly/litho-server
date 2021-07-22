package service

import (
	"context"
	"github.com/kainonly/gin-helper/hash"
	"lab-api/model"
)

type Admin struct {
	key string

	*Dependent
}

func NewAdmin(dep *Dependent) *Admin {
	return &Admin{
		Dependent: dep,
		key:       dep.Config.App.Key("sys:admin"),
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

func (x *Admin) Clear(ctx context.Context) error {
	return x.Redis.Del(ctx, x.key).Err()
}
