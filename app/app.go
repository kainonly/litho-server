package app

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
	"van-api/app/controller"
	"van-api/app/middleware/cors"
	"van-api/app/types"
)

func Application(option *types.Config) {
	app := iris.Default()
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.Cors(types.CorsOption{}))
	control := controller.New()
	app.Get("/", control.Default)
	app.Options("*", control.Option)
	main := app.Party("/main")
	{
		main.Post("/verify", control.MainVerify)
	}
	app.Listen(option.Listen)
}
