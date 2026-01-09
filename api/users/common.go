package users

import (
	"server/common"

	"go.uber.org/fx"
)

var Provides = fx.Provide(
	func(i *Service) *Controller {
		return &Controller{UsersX: i}
	},
	func(i common.Inject) *Service {
		return &Service{Inject: &i}
	},
)

type Controller struct {
	UsersX *Service
}

type Service struct {
	*common.Inject
}
