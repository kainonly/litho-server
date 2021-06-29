package main

import (
	"github.com/kainonly/gin-planx"
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
			bootstrap.HttpServer,
			planx.Initialize,
		),
		service.Provides,
		controller.Provides,
		fx.Invoke(routes.Initialize),
	).Run()
}
