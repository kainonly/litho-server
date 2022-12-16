package index

import (
	"github.com/google/wire"
)

var Provides = wire.NewSet(
	wire.Struct(new(Controller), "*"),
	wire.Struct(new(Service), "*"),
)

var MsgAuthenticationExpired = "authentication has expired, please log in again"
