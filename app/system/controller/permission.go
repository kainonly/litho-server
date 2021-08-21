package controller

import (
	"github.com/kainonly/go-bit/crud"
	"lab-api/model"
)

type Permission struct {
	*Dependency
	*crud.Resource
}

func NewPermission(d Dependency) *Permission {
	return &Permission{
		Dependency: &d,
		Resource:   d.Crud.Make(model.Resource{}),
	}
}
