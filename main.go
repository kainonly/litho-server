package main

import (
	curd "github.com/kainonly/gin-curd"
	"go.uber.org/fx"
	"lab-api/application"
	"lab-api/application/redis"
	"lab-api/bootstrap"
)

func main() {
	fx.New(
		fx.NopLogger,
		fx.Provide(
			bootstrap.LoadConfiguration,
			bootstrap.InitializeDatabase,
			bootstrap.InitializeRedis,
			bootstrap.HttpServer,
			redis.Initialize,
			curd.Initialize,
		),
		fx.Invoke(application.Application),
	).Run()
}
