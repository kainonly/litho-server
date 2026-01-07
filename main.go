package main

import (
	"server/api"
	"server/bootstrap"

	"go.uber.org/fx"
)

func main() {
	fx.New(
		bootstrap.Provides,
		api.Options,
	).Run()
}
