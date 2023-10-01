package index

import (
	"github.com/google/wire"
)

var Provides = wire.NewSet(
	wire.Struct(new(Controller), "*"),
	wire.Struct(new(Service), "*"),
)

type M = map[string]interface{}

func R(code string, msg string) M {
	return M{
		"code": code,
		"msg":  msg,
	}
}
