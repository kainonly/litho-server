package controller

import (
	"github.com/gin-gonic/gin"
)

type Index struct {
	*Services
}

func NewIndex(s Services) *Index {
	return &Index{&s}
}

func (x *Index) Index(c *gin.Context) interface{} {
	return gin.H{
		"version": "1.0",
	}
}
