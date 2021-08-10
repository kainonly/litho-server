package service

import (
	"lab-api/model"
)

type Admin struct {
	*Dependent
}

func NewAdmin(dep Dependent) *Admin {
	return &Admin{
		Dependent: &dep,
	}
}

func (x *Admin) FindOne(query Query) (data model.Admin, err error) {
	if err = query(x.Db).First(&data).Error; err != nil {
		return
	}
	return
}
