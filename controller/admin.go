package controller

import (
	"github.com/kainonly/go-bit/crud"
	"lab-api/model"
	"lab-api/service"
)

type Admin struct {
	*service.Services
	*crud.Resource
}

func NewAdmin(s *service.Services) *Admin {
	return &Admin{s, s.Crud.Make(model.Admin{})}
}
