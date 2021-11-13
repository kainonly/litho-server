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

func API() (*gin.Engine, error) {
	wire.Build(
		wire.Struct(new(common.App), "*"),
		bootstrap.Provides,
		app.Provides,
	)
	return &gin.Engine{}, nil
}
