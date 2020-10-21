package controller

import (
	"github.com/kataras/iris/v12/mvc"
	"van-api/app/model"
	"van-api/curd"
	"van-api/helper/res"
)

type AclController struct {
}

func (c *AclController) BeforeActivation(b mvc.BeforeActivation) {
}

func (c *AclController) PostOriginlists(curd *curd.Curd) interface{} {
	curd.Originlists(model.Acl{})
	return res.Ok()
}
