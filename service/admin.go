package service

import (
	"lab-api/model"
)

type Admin struct {
	*Dependency
}

func NewAdmin(dep Dependency) *Admin {
	return &Admin{&dep}
}

func (x *Admin) Create() {

}

func (x *Admin) FindAll() {

}

func (x *Admin) FindOne(query Query) (data model.Admin, err error) {
	if err = query(x.Db).First(&data).Error; err != nil {
		return
	}
	return
}

func (x *Admin) Update() {

}

func (x *Admin) Remove() {

}
