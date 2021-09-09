//go:build wireinject

package main

import (
	"github.com/google/wire"
	"lab-api/app"
	"lab-api/common"
)

func App(_ *common.Set) (*app.App, error) {
	wire.Build(
		common.HttpServer,
		common.InitializeDatabase,
		common.InitializeRedis,
		common.InitializeCookie,
		common.InitializeAuthx,
		common.InitializeCipher,
		wire.Struct(new(common.Dependency), "*"),
		app.Provides,
	)
	return &app.App{}, nil
}
