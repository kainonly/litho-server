package users

import (
	"context"
	"server/api/sessions"
	"server/common"

	"github.com/goforj/wire"
)

var Provides = wire.NewSet(
	wire.Struct(new(Controller), "*"),
	wire.Struct(new(Service), "*"),
)

type Controller struct {
	UsersX *Service
}

type Service struct {
	*common.Inject

	SessionsX *sessions.Service
}

func (x *Service) RefreshCache(ctx context.Context) error {
	return x.RDb.Del(ctx, "iam:users").Err()
}

func (x *Service) GetIAMUser(ctx context.Context, id string) (result *common.IAMUser, err error) {
	return
}
