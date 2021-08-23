package controller

import (
	"github.com/kainonly/go-bit/crud"
	"lab-api/model"
)

type Acl struct {
	*Dependency
	*crud.Crud
}

func NewAcl(d Dependency) *Acl {
	return &Acl{
		Dependency: &d,
		Crud:       crud.New(d.Db, model.Acl{}),
	}
}
