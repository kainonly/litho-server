package roles

import (
	"server/common"

	"go.uber.org/fx"
)

var Provides = fx.Provide(
	func(i *Service) *Controller { return &Controller{RolesX: i} },
	func(i common.Inject) *Service { return &Service{Inject: &i} },
)

type Controller struct {
	RolesX *Service
}

type Service struct {
	*common.Inject
}
