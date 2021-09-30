package main

import (
	"go.uber.org/fx"
	"laboratory/api"
	"laboratory/bootstrap"
)

func main() {
	fx.New(
		fx.NopLogger,
		bootstrap.Provides,
		api.Options,
	).Run()
}
