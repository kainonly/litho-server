package orgs

import (
	"server/common"

	"go.uber.org/fx"
)

var Provides = fx.Provide(
	func(i *Service) *Controller {
		return &Controller{OrgsX: i}
	},
	func(i common.Inject) *Service {
		return &Service{Inject: &i}
	},
)

type Controller struct {
	OrgsX *Service
}

type Service struct {
	*common.Inject
}
