// +build wireinject
// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/google/wire"
	"lab-api/controller"
	"lab-api/service"
)

func Boot() (*controller.Controllers, error) {
	wire.Build(
		LoadConfiguration,
		InitializeDatabase,
		InitializeRedis,
		service.Provides,
		controller.Provides,
	)
	return nil, nil
}
