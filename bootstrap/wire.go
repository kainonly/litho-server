//go:build wireinject
// +build wireinject

package bootstrap

import (
	"github.com/google/wire"
	"github.com/weplanx/server/api"
	"github.com/weplanx/server/common"
	"github.com/weplanx/server/utils/locker"
)

func NewAPI() (*api.API, error) {
	wire.Build(
		LoadStaticValues,
		UseGorm,
		UseRedis,
		UseNats,
		UseJetStream,
		UseStore,
		UseHertz,
		UseTransfer,
		api.Provides,
		wire.Struct(new(api.API), "*"),
		wire.Struct(new(common.Inject), "*"),
		wire.Struct(new(locker.Locker), "*"),
	)
	return &api.API{}, nil
}
