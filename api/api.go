package api

import (
	"context"
	"server/api/index"
	"server/common"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/google/wire"
	"github.com/weplanx/go/csrf"
)

var Provides = wire.NewSet(
	index.Provides,
)

type API struct {
	*common.Inject

	Hertz  *server.Hertz
	Csrf   *csrf.Csrf
	Index  *index.Controller
	IndexX *index.Service
}

func (x *API) SetupRoutes(h *server.Hertz) (err error) {
	//csrfToken := x.Csrf.VerifyToken(!x.V.IsRelease())
	//auth := x.AuthGuard()

	h.GET("", x.Index.Ping)

	//m := []app.HandlerFunc{csrfToken, auth, audit}
	return
}

func (x *API) AuthGuard() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		ts := c.Cookie("TOKEN")
		if ts == nil {
			c.AbortWithStatusJSON(401, utils.H{
				"code":    0,
				"message": "authentication has expired, please log in again",
			})
			return
		}

		c.Next(ctx)
	}
}

func (x *API) Initialize(ctx context.Context) (h *server.Hertz, err error) {
	h = x.Hertz

	return
}
