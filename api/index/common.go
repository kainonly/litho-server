package index

import (
	"github.com/weplanx/server/common"
	"go.uber.org/fx"
)

var Provides = fx.Provide(
	func(index *Service) *Controller {
		return &Controller{
			IndexService: index,
		}
	},
	func(i common.Inject) *Service {
		return &Service{
			Inject: &i,
		}
	},
)
