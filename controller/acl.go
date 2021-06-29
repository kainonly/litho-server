package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/kainonly/gin-helper/mvc"
)

type Acl struct{}

func NewAcl() *Acl {
	return &Acl{}
}

func (x *Acl) Index(c *gin.Context) interface{} {
	return mvc.Ok{
		"msg": "hello",
	}
}
