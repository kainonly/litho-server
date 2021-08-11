package controller

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"lab-api/service"
)

type Index struct {
	admin *service.Admin
}

func NewIndex(admin *service.Admin) *Index {
	return &Index{admin}
}

func (x *Index) Index(c *gin.Context) interface{} {
	data, err := x.admin.FindOne(func(tx *gorm.DB) *gorm.DB {
		return tx.
			Where("username = ?", "kain").
			Where("status = ?", true)
	})
	if err != nil {
		return err
	}
	return gin.H{
		"version": "1.0",
		"data":    data,
	}
}
