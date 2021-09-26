package system

import (
	"go.uber.org/fx"
)

var Provides = fx.Provide(
	func(i InjectController) *Controller {
		return &Controller{
			InjectController: &i,
		}
	},
	func(i InjectService) *Service {
		return &Service{
			InjectService:   &i,
			VerificationKey: i.Set.RedisKey("verify:"),
		}
	},
)
