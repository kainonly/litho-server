package service

import "github.com/gin-gonic/gin"

type Index struct {
	*Dependency
}

func NewIndex(d Dependency) *Index {
	return &Index{
		Dependency: &d,
	}
}

func (x *Index) Version() interface{} {
	return gin.H{
		"version": "1.1",
	}
}
