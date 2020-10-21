package app

import (
	"github.com/kataras/iris/v12/mvc"
	"van-api/app/controller"
)

func Application(app *mvc.Application) {
	app.Party("/main").Handle(new(controller.MainController))
	app.Party("/acl").Handle(new(controller.AclController))
}
