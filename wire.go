//go:build wireinject

package main

import (
	"github.com/google/wire"
	"lab-api/app"
	"lab-api/common"
)

func Bootstrap(_ *common.Set) (*app.App, error) {
	wire.Build(app.Provides)
	return &app.App{}, nil
}
