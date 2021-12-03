package app

import (
	"api/app/index"
	"github.com/gofiber/fiber/v2"
	"github.com/weplanx/go/api"
	"github.com/weplanx/go/passport"
	"github.com/weplanx/go/route"
	"go.uber.org/fx"
)

var Options = fx.Options(
	index.Provides,
	fx.Invoke(func(
		app *fiber.App,
		api *api.API,
		passport *passport.Passport,
		index *index.Controller,
	) {
		app.Get("/", route.Returns(Home))
		auth := AuthGuard(passport)
		route.Auto(app, index)
		route.Auto(app, api, route.SetPath(":collection"), route.SetMiddleware(auth))
	}),
)
