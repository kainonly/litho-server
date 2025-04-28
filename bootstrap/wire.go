//go:build wireinject
// +build wireinject

package bootstrap

import (
	"github.com/google/wire"
	"server/api"
	"server/common"
)

func NewAPI(values *common.Values) (*api.API, error) {
	wire.Build(
		wire.Struct(new(api.API), "*"),
		wire.Struct(new(common.Inject), "*"),
		UseGorm,
		UseRedis,
		UsePassport,
		UseCsrf,
		UseCipher,
		UseLocker,
		UseCaptcha,
		UseHertz,
		api.Provides,
	)
	return &api.API{}, nil
}
