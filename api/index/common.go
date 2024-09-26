package index

import (
	"github.com/google/wire"
	"github.com/weplanx/go/csrf"
	"github.com/weplanx/server/common"
)

var Provides = wire.NewSet(
	wire.Struct(new(Controller), "*"),
	wire.Struct(new(Service), "*"),
)

type Controller struct {
	V    *common.Values
	Csrf *csrf.Csrf

	IndexX *Service
}

type Service struct {
	*common.Inject

	Passport *common.APIPassport
}

type M = map[string]interface{}

func R(code string, msg string) M {
	return M{
		"code": code,
		"msg":  msg,
	}
}
