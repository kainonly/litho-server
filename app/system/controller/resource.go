package controller

import (
	"github.com/kainonly/go-bit/crud"
	"lab-api/model"
)

type Resource struct {
	*Dependency
	*crud.Resource
}

func NewResource(d Dependency) *Resource {
	return &Resource{
		Dependency: &d,
		Resource:   d.Crud.Make(model.Resource{}),
	}
}
