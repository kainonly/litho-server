package main

import (
	"github.com/kainonly/go-bit"
	"go.uber.org/fx"
	"lab-api/bootstrap"
	"lab-api/controller"
	"lab-api/routes"
	"lab-api/service"
)

func main() {
	fx.New(
		fx.NopLogger,
		fx.Provide(
			bit.LoadConfiguration,
			bit.InitializeCrud,
			bit.InitializeCookie,
			bootstrap.InitializeDatabase,
			bootstrap.InitializeRedis,
			bootstrap.HttpServer,
		),
		service.Provides,
		controller.Provides,
		fx.Invoke(routes.Initialize),
	).Run()
}
