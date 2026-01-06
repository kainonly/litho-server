package index

import (
	"server/common"

	"go.uber.org/fx"
)

type Controller struct {
	IndexX *Service
}

type Service struct {
	*common.Inject
}

var Provides = fx.Provide(
	func(i *Service) *Controller {
		return &Controller{IndexX: i}
	},
	func(i common.Inject) *Service {
		return &Service{Inject: &i}
	},
)
