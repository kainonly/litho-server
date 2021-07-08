package main

import (
	bit "github.com/kainonly/gin-bit"
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
			bootstrap.LoadConfiguration,
			bootstrap.InitializeDatabase,
			bootstrap.InitializeCookie,
			bootstrap.InitializeAuth,
			bootstrap.HttpServer,
			bit.Initialize,
		),
		service.Provides,
		controller.Provides,
		fx.Invoke(routes.Initialize),
	).Run()
}
