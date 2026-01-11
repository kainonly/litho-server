package menus

import (
	"server/common"

	"go.uber.org/fx"
)

var Provides = fx.Provide(
	func(i *Service) *Controller { return &Controller{MenusX: i} },
	func(i common.Inject) *Service { return &Service{Inject: &i} },
)

type Controller struct {
	MenusX *Service
}

type Service struct {
	*common.Inject
}
