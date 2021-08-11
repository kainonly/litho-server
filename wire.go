// +build wireinject
// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/google/wire"
	"github.com/kainonly/gin-bit"
	"lab-api/controller"
	"lab-api/service"
)

func Bootstrap() (*controller.Controllers, error) {
	wire.Build(
		gin_bit.LoadConfiguration,
		gin_bit.InitializeDatabase,
		gin_bit.InitializeRedis,
		service.Provides,
		controller.Provides,
	)
	return nil, nil
}
