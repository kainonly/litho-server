package service

import (
	"github.com/kainonly/gin-helper/hash"
	"github.com/kainonly/gin-helper/rbac"
	"lab-api/model"
)

type Admin struct {
	rbac.UserFn
	*Dependent

	key string
}

func NewAdmin(dep *Dependent) *Admin {
	return &Admin{
		Dependent: dep,
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
