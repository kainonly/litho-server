package bootstrap

import (
	"github.com/google/wire"
	"github.com/nats-io/nats.go"
	"github.com/weplanx/support"
	"github.com/weplanx/support/bootstrap"
	"github.com/weplanx/support/utils/captcha"
	"github.com/weplanx/support/utils/locker"
	"github.com/weplanx/transfer"
)

var Provides = wire.NewSet(
	bootstrap.LoadStaticValues,
	bootstrap.UseMongoDB,
	bootstrap.UseDatabase,
	bootstrap.UseRedis,
	bootstrap.UseNats,
	bootstrap.UseJetStream,
	bootstrap.UseHertz,
	UseTransfer,
	wire.Struct(new(captcha.Captcha), "*"),
	wire.Struct(new(locker.Locker), "*"),
)

// UseTransfer 初始日志传输
// https://github.com/weplanx/transfer
func UseTransfer(values *support.Values, js nats.JetStreamContext) (*transfer.Transfer, error) {
	return transfer.New(values.Namespace, js)
}
