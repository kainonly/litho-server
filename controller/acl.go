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

func (x *Acl) Get(c *gin.Context) interface{} {
	var body struct {
		bit.Conditions `json:"where" binding:"required_without=Id,gte=0,dive,len=3,dive,required"`
	}
	bit.Complex(c,
		bit.SetBody(&body),
	)
	return x.Crud.Get(c)
}
