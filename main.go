package main

import (
	"go.uber.org/fx"
	"lab-api/bootstrap"
)

func main() {
	fx.New(
		fx.NopLogger,
		fx.Provide(
			bootstrap.LoadConfiguration,
			bootstrap.InitializeDatabase,
			bootstrap.HttpServer,
		),
	).Run()
}
