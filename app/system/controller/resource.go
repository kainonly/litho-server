package controller

import (
	"github.com/kainonly/go-bit/crud"
	"lab-api/model"
)

type Resource struct {
	*Dependency
	*crud.Crud
}

func NewResource(d Dependency) *Resource {
	return &Resource{
		Dependency: &d,
		Crud:       crud.New(d.Db, model.Resource{}),
	}
}
