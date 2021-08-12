package controller

import (
	"github.com/kainonly/go-bit/crud"
	"lab-api/model"
)

type Resource struct {
	*Services
	*crud.Resource
}

func NewResource(s Services) *Resource {
	return &Resource{&s, s.Crud.Make(model.Resource{})}
}
