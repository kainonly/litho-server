// +build wireinject
// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/google/wire"
	"github.com/kainonly/go-bit"
	"lab-api/controller"
	"lab-api/service"
)

func Boot(config bit.Config) (*controller.Controllers, error) {
	wire.Build(
		InitializeDatabase,
		InitializeRedis,
		bit.Initialize,
		service.Provides,
		controller.Provides,
	)
	return nil, nil
}
