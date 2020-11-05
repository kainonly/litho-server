package main

import (
	curd "github.com/kainonly/gin-curd"
	"go.uber.org/fx"
	"taste-api/application"
	"taste-api/application/cache"
	"taste-api/bootstrap"
)

func main() {
	fx.New(
		fx.Provide(
			bootstrap.LoadConfiguration,
			bootstrap.InitializeDatabase,
			bootstrap.InitializeRedis,
			bootstrap.HttpServer,
			cache.Initialize,
			curd.Initialize,
		),
		fx.Invoke(application.Application),
	).Run()
}
