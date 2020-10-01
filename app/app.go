package app

import (
	"github.com/kataras/iris/v12"
	"van-api/app/controller"
	"van-api/app/middleware/cors"
	"van-api/app/types"
)

func Application(option *types.Config) {
	app := iris.Default()
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
