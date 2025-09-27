package index

import (
	"server/common"

	"github.com/google/wire"
	"github.com/weplanx/go/passport"
)

var Provides = wire.NewSet(
	wire.Struct(new(Controller), "*"),
	wire.Struct(new(Service), "*"),
)

type Controller struct {
	V *common.Values

	IndexX *Service
}

type Service struct {
	*common.Inject

	Passport *passport.Passport
}

type M = map[string]interface{}

func R(code string, msg string) M {
	return M{
		"code": code,
		"msg":  msg,
	}
}
