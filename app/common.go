package app

import (
	"api/app/index"
	"api/app/pages"
	"api/app/roles"
	"github.com/gofiber/fiber/v2"
	"github.com/weplanx/go/api"
	"github.com/weplanx/go/passport"
	"github.com/weplanx/go/route"
	"go.uber.org/fx"
)

var Options = fx.Options(
	index.Provides,
	pages.Provides,
	roles.Provides,
	fx.Provide(api.New, api.AutoController),
	fx.Invoke(func(
		app *fiber.App,
		api *api.Controller,
		passport *passport.Passport,
		index *index.Controller,
		pages *pages.Controller,
		roles *roles.Controller,
	) {
		app.Get("/", route.Returns(Home))

		//auth := AuthGuard(passport)
		route.Auto(app, index)
		route.Auto(app, pages, route.SetPath("/pages"))
		route.Auto(app, roles, route.SetPath("/roles"))
		route.Auto(app, api, route.SetPath(":collection"))
		//route.Auto(app, index, route.SetMiddleware(auth, "Code", "RefreshToken", "Logout", "Pages"))
		//route.Auto(app, api, route.SetPath(":collection"), route.SetMiddleware(auth))
	}),
)
