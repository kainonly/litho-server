package app

import (
	"github.com/kataras/iris/v12/mvc"
	"van-api/app/controller"
	"van-api/app/controller/acl"
)

func Application(app *mvc.Application) {
	app.Party("/main").Handle(new(controller.Controller))
	app.Party("/acl").Handle(new(acl.Controller))
}
