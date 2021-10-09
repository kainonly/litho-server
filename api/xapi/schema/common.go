package schema

import (
	"github.com/weplanx/support/api"
	"go.uber.org/fx"
)

var Provides = fx.Provide(
	func(i InjectController) *Controller {
		return &Controller{
			InjectController: &i,
			API:              api.New(i.Mongo, i.Db, api.SetCollection("schema")),
		}
	},
	func(i InjectService) *Service {
		return &Service{
			InjectService: &i,
		}
	},
)
