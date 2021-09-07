package controller

import (
	"github.com/gin-gonic/gin"
)

type Index struct {
	*Dependency
}

func NewIndex(d Dependency) *Index {
	return &Index{
		Dependency: &d,
	}
}

func (x *Index) Index(c *gin.Context) interface{} {
	return x.IndexService.Version()
}
