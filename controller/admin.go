package controller

import (
	"github.com/kainonly/go-bit/crud"
	"lab-api/model"
)

type Admin struct {
	*Services
	*crud.Resource
}

func NewAdmin(s Services) *Admin {
	return &Admin{&s, s.Crud.Make(model.Admin{})}
}
