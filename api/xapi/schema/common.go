package schema

import (
	"go.uber.org/fx"
)

var Provides = fx.Provide(
	func(i InjectController) *Controller {
		return &Controller{
			InjectController: &i,
			//API:              api.New(i.Db, api.SetModel("schema")),
		}
	},
	func(i InjectService) *Service {
		return &Service{
			InjectService: &i,
		}
	},
)
