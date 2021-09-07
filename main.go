package main

import (
	"go.uber.org/fx"
	"lab-api/app/api"
	"lab-api/app/system"
	"lab-api/bootstrap"
)

func main() {
	fx.New(
		fx.NopLogger,
		fx.Provide(
			bootstrap.LoadConfiguration,
			bootstrap.InitializeDatabase,
			bootstrap.InitializeRedis,
			bootstrap.InitializeCookie,
			bootstrap.InitializeAuthx,
			bootstrap.InitializeCipher,
			bootstrap.HttpServer,
		),
		system.App,
		api.App,
	).Run()
}
