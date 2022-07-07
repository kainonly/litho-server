//go:build wireinject
// +build wireinject

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"server/app"
	"server/bootstrap"
	"server/common"
)

func App(value *common.Values) (*gin.Engine, error) {
	wire.Build(
		wire.Struct(new(common.Inject), "*"),
		bootstrap.Provides,
		app.Provides,
	)
	return &gin.Engine{}, nil
}
