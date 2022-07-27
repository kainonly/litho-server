//go:build wireinject
// +build wireinject

package bootstrap

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/google/wire"
	"github.com/weplanx/server/api"
	"github.com/weplanx/server/common"
)

func OkLetsGo(value *common.Values) (*server.Hertz, error) {
	wire.Build(
		wire.Struct(new(common.Inject), "*"),
		Provides,
		api.Provides,
	)
	return &server.Hertz{}, nil
}
