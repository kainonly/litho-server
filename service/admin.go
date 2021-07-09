package service

import (
	"github.com/kainonly/gin-helper/hash"
	"gorm.io/gorm"
	"lab-api/model"
)

type Admin struct {
	db *gorm.DB
}

func NewAdmin(db *gorm.DB) *Admin {
	return &Admin{
		db: db,
	}
}

func (x *Admin) FindOne(query Query) (data model.Admin, err error) {
	if err = query(x.db).First(&data).Error; err != nil {
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
