//go:build wireinject
// +build wireinject

package main

import (
	"api/app"
	"api/bootstrap"
	"api/common"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func App(value *common.Values) (*gin.Engine, error) {
	wire.Build(
		wire.Struct(new(common.Inject), "*"),
		bootstrap.Provides,
		app.Provides,
	)
	return &gin.Engine{}, nil
}
