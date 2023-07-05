package main

import (
	"github.com/weplanx/server/api"
	"github.com/weplanx/server/bootstrap"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		fx.NopLogger,
		fx.Provide(
			bootstrap.LoadStaticValues,
			bootstrap.UseMongoDB,
			bootstrap.UseRedis,
			bootstrap.UseHertz,
		),
		api.Options,
	).Run()
}
