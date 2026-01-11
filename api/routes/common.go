package routes

import (
	"server/common"

	"go.uber.org/fx"
)

var Provides = fx.Provide(
	func(i *Service) *Controller { return &Controller{RoutesX: i} },
	func(i common.Inject) *Service { return &Service{Inject: &i} },
)

type Controller struct {
	RoutesX *Service
}
type Service struct {
	*common.Inject
}
