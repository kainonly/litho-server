//go:build wireinject
// +build wireinject

package bootstrap

import (
	"github.com/google/wire"
	"github.com/weplanx/server/api"
	"github.com/weplanx/server/common"
	"github.com/weplanx/server/openapi"
	"github.com/weplanx/server/xapi"
)

func NewAPI(values *common.Values) (*api.API, error) {
	wire.Build(
		wire.Struct(new(api.API), "*"),
		wire.Struct(new(common.Inject), "*"),
		UseMongoDB,
		UseDatabase,
		UseRedis,
		UseNats,
		UseJetStream,
		UseKeyValue,
		UseValues,
		UseSessions,
		UseRest,
		UseCsrf,
		UseCipher,
		UseAPIPassport,
		UseLocker,
		UseCaptcha,
		UseTransfer,
		UseInflux,
		UseHertz,
		api.Provides,
	)
	return &api.API{}, nil
}

func NewXAPI(values *common.Values) (*xapi.API, error) {
	wire.Build(
		wire.Struct(new(xapi.API), "*"),
		wire.Struct(new(common.Inject), "*"),
		UseMongoDB,
		UseDatabase,
		UseRedis,
		UseNats,
		UseJetStream,
		UseKeyValue,
		UseCipher,
		UseLocker,
		UseCaptcha,
		UseTransfer,
		UseHertz,
		xapi.Provides,
	)
	return &xapi.API{}, nil
}

func NewOpenAPI(values *common.Values) (*openapi.API, error) {
	wire.Build(
		wire.Struct(new(openapi.API), "*"),
		wire.Struct(new(common.Inject), "*"),
		UseMongoDB,
		UseDatabase,
		UseRedis,
		UseNats,
		UseJetStream,
		UseKeyValue,
		UseCipher,
		UseLocker,
		UseCaptcha,
		UseTransfer,
		UseHertz,
		openapi.Provides,
	)
	return &openapi.API{}, nil
}
