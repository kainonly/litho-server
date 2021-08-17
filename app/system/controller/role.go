package controller

import (
	"github.com/kainonly/go-bit/crud"
	"lab-api/model"
)

type Role struct {
	*Dependency
	*crud.Resource
}

func NewRole(d Dependency) *Role {
	return &Role{
		Dependency: &d,
		Resource:   d.Crud.Make(model.Role{}),
	}
}
