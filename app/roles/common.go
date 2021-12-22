package roles

import (
	"github.com/weplanx/go/api"
	"go.uber.org/fx"
)

var Provides = fx.Provide(
	func(i InjectController) *Controller {
		return &Controller{
			Controller:       api.SetController(i.APIx, "roles"),
			InjectController: &i,
		}
	},
	func(i InjectService) *Service {
		return &Service{
			InjectService: &i,
		}
	},
)
