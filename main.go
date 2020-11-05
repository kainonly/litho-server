package main

import (
	"go.uber.org/fx"
	"taste-api/application"
	"taste-api/bootstrap"
)

func main() {
	fx.New(
		fx.Provide(
			bootstrap.LoadConfiguration,
			bootstrap.InitializeDatabase,
			bootstrap.InitializeRedis,
			bootstrap.HttpServer,
		),
		fx.Invoke(application.Application),
	).Run()
}
