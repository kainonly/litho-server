package controller

import (
	"github.com/gin-gonic/gin"
	"lab-api/service"
)

type Index struct {
	*service.Services
}

func NewIndex(s *service.Services) *Index {
	return &Index{s}
}

func (x *Index) Index(c *gin.Context) interface{} {
	return gin.H{
		"version": "1.0",
	}
}
