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
		bootstrap.Provides,
		api.Options,
		system.Options,
	).Run()
}
