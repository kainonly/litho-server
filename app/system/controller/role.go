package controller

import (
	"github.com/kainonly/go-bit/crud"
	"lab-api/model"
)

type Role struct {
	*Dependency
	*crud.Crud
}

func NewRole(d *Dependency) *Role {
	return &Role{
		Dependency: d,
		Crud:       crud.New(d.Db, &model.Role{}),
	}
}
