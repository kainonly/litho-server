package main

import (
	"api/app"
	"api/bootstrap"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		fx.NopLogger,
		bootstrap.Provides,
		app.Options,
	).Run()
}
