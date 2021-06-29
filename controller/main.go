package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/kainonly/gin-planx/mvc"
)

type Main struct{}

func NewMain() *Main {
	return &Main{}
}

func (x *Main) Index(c *gin.Context) interface{} {
	return mvc.Ok{
		"msg": "hello",
	}
}
