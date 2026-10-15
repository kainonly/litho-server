package orgs

import (
	"context"
	"litho-api/common"

	"github.com/goforj/wire"
	"github.com/kainonly/go/help"
)

var Provides = wire.NewSet(
	wire.Struct(new(Controller), "*"),
	wire.Struct(new(Service), "*"),
)

const Key = "orgs"

type Controller struct {
	V *common.Values

	OrgsX *Service
}

type Service struct {
	*common.Inject
}

type M = map[string]any

var ErrAccess = help.E(0, `您不具备系统权限，请联系超级管理员进行管理`)

func (x *Service) CheckAccess(ctx context.Context, user *common.IAMUser, ids ...string) (err error) {
	if err = user.Can(`WM`); err != nil {
		return ErrAccess
	}
	return
}
