package controller

import (
	"github.com/kainonly/go-bit/crud"
	"lab-api/model"
)

type Permission struct {
	*Dependency
	*crud.Crud
}

func NewPermission(d Dependency) *Permission {
	return &Permission{
		Dependency: &d,
		Crud:       crud.New(d.Db, model.Resource{}),
	}
}
