package api

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/google/wire"
	"github.com/hertz-contrib/jwt"
	"github.com/weplanx/api/api/pages"
	"github.com/weplanx/server/api"
)

var Provides = wire.NewSet(
	wire.Struct(new(API), "*"),
	pages.Provides,
)

type API struct {
	*api.API

	PagesController *pages.Controller
	PagesService    *pages.Service
}

func (x *API) Routes(h *server.Hertz) (auth *jwt.HertzJWTMiddleware, err error) {
	if auth, err = x.API.Routes(h); err != nil {
		return
	}

	_pages := h.Group("pages", auth.MiddlewareFunc())
	{
		_pages.GET(":id", x.PagesController.GetOne)
		_pages.GET(":id/indexes", x.PagesController.GetIndexes)
		_pages.PUT(":id/indexes/:index", x.PagesController.SetIndex)
		_pages.DELETE(":id/indexes/:index", x.PagesController.DeleteIndex)
	}

	return
}
