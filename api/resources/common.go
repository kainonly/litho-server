package resources

import (
	"server/common"

	"go.uber.org/fx"
)

var Provides = fx.Provide(
	func(i *Service) *Controller { return &Controller{ResourcesX: i} },
	func(i common.Inject) *Service { return &Service{Inject: &i} },
)

type Controller struct {
	ResourcesX *Service
}

type Service struct {
	*common.Inject
}
