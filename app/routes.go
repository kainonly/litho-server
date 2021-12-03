package app

import (
	"api/app/index"
	"github.com/gofiber/fiber/v2"
	wpx "github.com/weplanx/go"
	"github.com/weplanx/go/api"
	"go.uber.org/fx"
)

var Options = fx.Options(
	index.Provides,
	fx.Invoke(func(
		app *fiber.App,
		api *api.API,
		index *index.Controller,
	) {
		app.Get("/", wpx.Returns(index.Index))
		app.Get("/get", wpx.Returns(index.Get))
		wpx.Auto(app, api, wpx.SetPath(":collection"))
	}),
)
