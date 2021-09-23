package main

import (
	"go.uber.org/fx"
	"lab-api/api"
	"lab-api/bootstrap"
)

func main() {
	fx.New(
		//fx.NopLogger,
		bootstrap.Provides,
		api.Options,
	).Run()
}
