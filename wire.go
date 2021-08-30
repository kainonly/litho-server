//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"lab-api/app"
	"lab-api/common"
)

func App(config common.Config) (*app.App, error) {
	wire.Build(app.Provides)
	return &app.App{}, nil
}
