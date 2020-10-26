package controller

import (
	"github.com/kataras/iris/v12/mvc"
	"gorm.io/gorm"
	"van-api/app/model"
	"van-api/curd"
)

type AclController struct {
}

func (c *AclController) BeforeActivation(b mvc.BeforeActivation) {
}

type TestBody struct {
	curd.Body
}

func (c *AclController) PostOriginlists(body *TestBody, mode *curd.Curd) interface{} {
	return mode.
		Originlists(model.Acl{}, body).
		Where(curd.Conditions{
			[]interface{}{"status", "=", "1"},
		}).
		Field([]string{"id", "name", "read", "write"}).
		Result()
}

func (c *AclController) PostTest(body *TestBody, mode *curd.Curd) interface{} {
	return mode.
		Get(model.Acl{}, body).
		Query(func(tx *gorm.DB) {
			tx.Where("id = ?", "2")
		}).
		Field([]string{"id", "name", "read", "write"}).
		Result()
}
