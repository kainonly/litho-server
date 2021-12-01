package x

import (
	"api/app/x/admin"
	"api/app/x/devops"
	"api/app/x/page"
	"api/app/x/schema"
	"go.uber.org/fx"
)

var Provides = fx.Provide(
	func(i InjectController) *Controller {
		return &Controller{
			InjectController: &i,
			Auth:             i.Passport.Make("system"),
		}
	},
	func(i InjectService) *Service {
		return &Service{
			InjectService: &i,
		}
	},
)

var Options = fx.Options(
	Provides,
	page.Provides,
	admin.Provides,
	devops.Provides,
	schema.Provides,
)
