//go:build wireinject
// +build wireinject

package test

import (
	"github.com/google/wire"
	"server/bootstrap"
	"server/common"
)

func Injectable(value *common.Values) (*common.Inject, error) {
	wire.Build(
		wire.Struct(new(common.Inject), "*"),
		bootstrap.Provides,
	)
	return &common.Inject{}, nil
}
