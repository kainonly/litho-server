package service

import (
	"lab-api/model"
)

type Admin struct {
	*Dependency
	Key string
}

func NewAdmin(d Dependency) *Admin {
	return &Admin{
		Dependency: &d,
		Key:        d.Config.RedisKey("admin"),
	}
}

func (x *Admin) FindByUsername(username string) (data model.Admin, err error) {
	if err = x.Db.
		Where("username = ?", username).
		Where("status = ?", true).
		First(&data).Error; err != nil {
		return
	}
	return
}
