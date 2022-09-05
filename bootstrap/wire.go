//go:build wireinject
// +build wireinject

package bootstrap

import (
	"github.com/google/wire"
	"github.com/weplanx/server/api"
	"github.com/weplanx/server/utils/captcha"
	"github.com/weplanx/server/utils/locker"
)

func NewAPI() (*api.API, error) {
	wire.Build(
		wire.Struct(new(api.API), "*"),
		LoadStaticValues,
		UseMongoDB,
		UseDatabase,
		UseRedis,
		UseNats,
		UseJetStream,
		UseHertz,
		UseTransfer,
		wire.Struct(new(captcha.Captcha), "*"),
		wire.Struct(new(locker.Locker), "*"),
		api.Provides,
	)
	return &api.API{}, nil
}
