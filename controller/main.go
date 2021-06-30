package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/kainonly/gin-helper/hash"
)

type Main struct{}

func NewMain() *Main {
	return &Main{}
}

func (x *Main) Index(c *gin.Context) interface{} {
	h, _ := hash.Make("hello")
	return gin.H{
		"val":  "hello",
		"hash": h,
	}
}
