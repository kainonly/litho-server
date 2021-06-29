package controller

import (
	bit "github.com/kainonly/gin-bit"
	"lab-api/model"
)

type Acl struct {
	*bit.Crud
}

func NewAcl(b *bit.Bit) *Acl {
	return &Acl{
		Crud: b.Crud(model.Acl{}),
	}
}
