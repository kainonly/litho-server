package controller

import (
	"github.com/gin-gonic/gin"
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

func (x *Acl) Add(c *gin.Context) interface{} {
	return x.Crud.Add(c)
}
