package controller

import (
	"github.com/kainonly/go-bit/crud"
	"lab-api/model"
)

type Acl struct {
	*Dependency
	*crud.Resource
}

func NewAcl(d Dependency) *Acl {
	return &Acl{
		Dependency: &d,
		Resource:   d.Crud.Make(model.Acl{}),
	}
}
