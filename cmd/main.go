package main

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
	"van-api/controller"
)

func main() {
	app := iris.Default()
	app.Use(recover.New())
	app.Use(logger.New())
	control := controller.New()
	route := app.Party("/")
	{
		route.Get("/", control.Index)
	}
	app.Listen(":8080")
}
