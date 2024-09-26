//go:build wireinject
// +build wireinject

package bootstrap

import (
	"github.com/google/wire"
	"github.com/weplanx/server/api"
	"github.com/weplanx/server/common"
)

func NewAPI(values *common.Values) (*api.API, error) {
	wire.Build(
		wire.Struct(new(api.API), "*"),
		wire.Struct(new(common.Inject), "*"),
		UseRedis,
		UseNats,
		UseJetStream,
		UseKeyValue,
		UseCsrf,
		UseCipher,
		UseAPIPassport,
		UseLocker,
		UseCaptcha,
		UseHertz,
		api.Provides,
	)
	return &api.API{}, nil
}
