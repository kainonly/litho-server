package service

import (
	"github.com/kainonly/gin-helper/hash"
	"gorm.io/gorm"
	"lab-api/model"
)

type Admins struct {
	db *gorm.DB
}

func NewAdmins(db *gorm.DB) *Admins {
	return &Admins{
		db: db,
	}
}

func (x *Admins) FindOne(query Query) (data model.Admin, err error) {
	if err = query(x.db).First(&data).Error; err != nil {
		return
	}
	return
}

func (x *Admins) Data(admin model.Admin) model.Admin {
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
