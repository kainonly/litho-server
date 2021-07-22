package controller

import "github.com/gin-gonic/gin"

type Main struct {
}

func NewMain() *Main {
	return &Main{}
}

func (x *Main) Index(c *gin.Context) interface{} {
	return gin.H{
		"version": "1.0",
	}
}
