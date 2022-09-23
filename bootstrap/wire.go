//go:build wireinject
// +build wireinject

package bootstrap

import (
	"github.com/google/wire"
	"github.com/weplanx/server/api"
	"github.com/weplanx/server/common"
	"github.com/weplanx/server/utils/locker"
)

func NewAPI(values *common.Values) (*api.API, error) {
	wire.Build(
		UseGorm,
		UseRedis,
		UseNats,
		UseJetStream,
		UseKeyValue,
		UseHertz,
		UseTransfer,
		api.Provides,
		wire.Struct(new(api.API), "*"),
		wire.Struct(new(common.Inject), "*"),
		wire.Struct(new(locker.Locker), "*"),
	)
	return &api.API{}, nil
}
