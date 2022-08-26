//go:build wireinject
// +build wireinject

package bootstrap

import (
	"github.com/google/wire"
	"github.com/weplanx/api/api"
	server "github.com/weplanx/server/api"
)

func NewAPI() (*api.API, error) {
	wire.Build(
		wire.Struct(new(server.API), "*"),
		server.Provides,
		Provides,
		api.Provides,
	)
	return &api.API{}, nil
}
