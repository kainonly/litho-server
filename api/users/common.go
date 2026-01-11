package users

import (
	"context"
	"server/api/sessions"
	"server/common"

	"go.uber.org/fx"
)

var Provides = fx.Provide(
	func(i *Service) *Controller { return &Controller{UsersX: i} },
	func(i common.Inject, s *sessions.Service) *Service { return &Service{Inject: &i, SessionsX: s} },
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
