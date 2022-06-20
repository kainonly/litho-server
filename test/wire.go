//go:build wireinject
// +build wireinject

package test

import (
	"api/bootstrap"
	"api/common"
	"github.com/google/wire"
)

func Injectable(value *common.Values) (*common.Inject, error) {
	wire.Build(
		wire.Struct(new(common.Inject), "*"),
		bootstrap.Provides,
	)
	return &common.Inject{}, nil
}
