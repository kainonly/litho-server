package main

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
	"van-api/controller"
	"van-api/middleware/cors"
)

func main() {
	app := iris.Default()
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.Cors(cors.Option{}))
	control := controller.New()
	app.Get("/", control.Default)
	app.Options("*", control.Option)
	main := app.Party("/main")
	{
		main.Post("/verify", control.MainVerify)
	}
	app.Listen(":8080")
}
