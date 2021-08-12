package controller

import (
	"github.com/kainonly/go-bit/crud"
	"lab-api/model"
	"lab-api/service"
)

type Resource struct {
	*service.Services
	*crud.Resource
}

func NewResource(s *service.Services) *Resource {
	return &Resource{s, s.Crud.Make(model.Resource{})}
}
