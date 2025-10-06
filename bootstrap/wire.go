//go:build wireinject
// +build wireinject

package bootstrap

import (
	"server/api"
	"server/common"

	"github.com/google/wire"
)

func NewAPI(values *common.Values) (*api.API, error) {
	wire.Build(
		wire.Struct(new(api.API), "*"),
		wire.Struct(new(common.Inject), "*"),
		UseGorm,
		UseRedis,
		UsePassport,
		UseCsrf,
		UseLocker,
		UseCaptcha,
		UseHertz,
		api.Provides,
	)
	return &api.API{}, nil
}
