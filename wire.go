//go:build wireinject
// +build wireinject

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/weplanx/server/api"
	"github.com/weplanx/server/bootstrap"
	"github.com/weplanx/server/common"
)

func OkLetsGo(value *common.Values) (*gin.Engine, error) {
	wire.Build(
		bootstrap.Provides,
		api.Provides,
	)
	return &gin.Engine{}, nil
}
