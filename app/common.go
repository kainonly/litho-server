package app

import (
	"api/app/admin"
	"api/app/index"
	"api/app/page"
	"github.com/gofiber/fiber/v2"
	"github.com/weplanx/go/api"
	"github.com/weplanx/go/passport"
	"github.com/weplanx/go/route"
	"go.uber.org/fx"
)

var Options = fx.Options(
	index.Provides,
	admin.Provides,
	page.Provides,
	fx.Invoke(func(
		app *fiber.App,
		api *api.API,
		passport *passport.Passport,
		index *index.Controller,
	) {
		auth := AuthGuard(passport)

		app.Get("/", route.Returns(Home))

		route.Auto(app, index, route.SetMiddleware(auth, "Code", "RefreshToken", "Logout", "Pages"))
		route.Auto(app, api, route.SetPath(":collection"), route.SetMiddleware(auth))
	}),
)
