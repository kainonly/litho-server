//go:build wireinject
// +build wireinject

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"laboratory/api"
	"laboratory/bootstrap"
	"laboratory/common"
)

func API() (*gin.Engine, error) {
	wire.Build(
		wire.Struct(new(common.App), "*"),
		bootstrap.Provides,
		api.Provides,
	)
	return &gin.Engine{}, nil
}
