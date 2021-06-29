package main

import (
	"go.uber.org/fx"
	"lab-api/bootstrap"
	"lab-api/controller"
	"lab-api/routes"
	"lab-api/service"
)

func main() {
	fx.New(
		fx.Provide(
			bootstrap.LoadConfiguration,
			bootstrap.InitializeDatabase,
			bootstrap.HttpServer,
		),
		service.Provides,
		controller.Provides,
		fx.Invoke(routes.Initialize),
	).Run()
}
