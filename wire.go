// +build wireinject
// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/google/wire"
	"github.com/kainonly/go-bit"
	"lab-api/bootstrap"
	"lab-api/controller"
	"lab-api/service"
)

func Boot(config bit.Config) (*controller.Controllers, error) {
	wire.Build(
		bootstrap.InitializeDatabase,
		bootstrap.InitializeRedis,
		bit.InitializeCrud,
		bit.InitializeCookie,
		service.Provides,
		controller.Provides,
	)
	return nil, nil
}
