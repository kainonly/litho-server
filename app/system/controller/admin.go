package controller

import (
	"github.com/kainonly/go-bit/crud"
	"lab-api/model"
)

type Admin struct {
	*Dependency
	*crud.Resource
}

func NewAdmin(d Dependency) *Admin {
	return &Admin{&d, d.Crud.Make(model.Admin{})}
}
