package app

import (
	"github.com/kataras/iris/v12/mvc"
	"van-api/app/controller"
	"van-api/app/controller/acl"
	"van-api/app/controller/admin"
	"van-api/app/controller/policy"
	"van-api/app/controller/resource"
	"van-api/app/controller/role"
)

func Application(app *mvc.Application) {
	app.Party("/main").Handle(new(controller.Controller))
	app.Party("/acl").Handle(new(acl.Controller))
	app.Party("/admin").Handle(new(admin.Controller))
	app.Party("/policy").Handle(new(policy.Controller))
	app.Party("/resource").Handle(new(resource.Controller))
	app.Party("/role").Handle(new(role.Controller))
}
