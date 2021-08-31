package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/kainonly/go-bit/crud"
	"gorm.io/gorm"
	"lab-api/model"
)

type Resource struct {
	*Dependency
	*crud.Crud
}

func NewResource(d *Dependency) *Resource {
	return &Resource{
		Dependency: d,
		Crud:       crud.New(d.Db, &model.Resource{}),
	}
}

func (x *Resource) OriginLists(c *gin.Context) interface{} {
	crud.Mix(c,
		crud.Query(func(tx *gorm.DB) *gorm.DB {
			return tx.Order("sort")
		}),
	)
	return x.Crud.OriginLists(c)
}
