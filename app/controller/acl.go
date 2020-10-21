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

func (c *AclController) PostOriginlists(curd *curd.Curd) interface{} {
	return curd.Originlists(model.Acl{})
}
