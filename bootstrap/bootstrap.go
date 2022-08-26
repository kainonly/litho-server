package bootstrap

import (
	"github.com/google/wire"
	"github.com/weplanx/server/bootstrap"
	"github.com/weplanx/server/utils/captcha"
	"github.com/weplanx/server/utils/locker"
)

var Provides = wire.NewSet(
	bootstrap.LoadStaticValues,
	bootstrap.UseMongoDB,
	bootstrap.UseDatabase,
	bootstrap.UseRedis,
	bootstrap.UseNats,
	bootstrap.UseJetStream,
	bootstrap.UseHertz,
	bootstrap.UseTransfer,
	wire.Struct(new(captcha.Captcha), "*"),
	wire.Struct(new(locker.Locker), "*"),
)
