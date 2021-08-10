// +build wireinject
// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/google/wire"
	"lab-api/bootstrap"
	"lab-api/controller"
	"lab-api/service"
)

func Bootstrap() (*controller.Controllers, error) {
	wire.Build(
		bootstrap.Provides,
		service.Provides,
		controller.Provides,
	)
	return nil, nil
}
