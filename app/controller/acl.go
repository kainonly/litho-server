package controller

import (
	"github.com/kataras/iris/v12/mvc"
	"van-api/app/model"
	"van-api/curd"
)

type AclController struct {
}

func (c *AclController) BeforeActivation(b mvc.BeforeActivation) {
}

func (c *AclController) PostOriginlists(body *curd.OriginListsBody, schema *curd.Curd) interface{} {
	return schema.Originlists(model.Acl{}, body).
		Where(curd.ArrayCondition{
			[]interface{}{"id=?", "1"},
		}).
		Field([]string{"id", "name"}).
		Result()
}
