package controller

import (
	"github.com/kainonly/go-bit/crud"
	"lab-api/model"
)

type Admin struct {
	*Dependency
	*crud.Crud
}

func NewAdmin(d *Dependency) *Admin {
	return &Admin{
		Dependency: d,
		Crud:       crud.New(d.Db, &model.Admin{}),
	}
}
