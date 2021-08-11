package controller

import (
	"github.com/kainonly/go-bit"
	"github.com/kainonly/go-bit/crud"
	"lab-api/model"
	"lab-api/service"
)

type Admin struct {
	admin *service.Admin
	crud  *crud.Crud
}

func NewAdmin(admin *service.Admin, bit *bit.Bit) *Admin {
	return &Admin{admin, bit.Crud(&model.Admin{})}
}
