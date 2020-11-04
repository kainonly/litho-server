package main

import (
	"go.uber.org/fx"
	"taste-api/bootstrap"
)

func main() {
	fx.New(
		fx.Provide(
			bootstrap.LoadConfiguration,
			bootstrap.InitializeDatabase,
			bootstrap.InitializeRedis,
		),
		fx.Invoke(bootstrap.HttpServer),
	).Run()
}
