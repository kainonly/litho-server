package app

import (
	"api/app/index"
	"api/app/pages"
	"github.com/gofiber/fiber/v2"
	"github.com/weplanx/go/api"
	"github.com/weplanx/go/passport"
	"github.com/weplanx/go/route"
	"go.uber.org/fx"
)

var Options = fx.Options(
	index.Provides,
	pages.Provides,
	fx.Invoke(func(
		app *fiber.App,
		api *api.API,
		passport *passport.Passport,
		index *index.Controller,
		pages *pages.Controller,
	) {
		app.Get("/", route.Returns(Home))

		//auth := AuthGuard(passport)
		route.Auto(app, index)
		route.Auto(app, pages, route.SetPath("/pages"))
		route.Auto(app, api, route.SetPath(":collection"))
		//route.Auto(app, index, route.SetMiddleware(auth, "Code", "RefreshToken", "Logout", "Pages"))
		//route.Auto(app, api, route.SetPath(":collection"), route.SetMiddleware(auth))
	}),
)
