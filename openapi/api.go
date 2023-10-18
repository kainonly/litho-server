package openapi

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/google/wire"
	"github.com/weplanx/server/common"
	"github.com/weplanx/server/openapi/index"
)

var Provides = wire.NewSet(
	index.Provides,
)

type API struct {
	*common.Inject

	Hertz        *server.Hertz
	Index        *index.Controller
	IndexService *index.Service
}

func (x *API) Routes(h *server.Hertz) (err error) {
	return
}

func (x *API) AuthGuard() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		ts := c.GetHeader("X-TOKEN")
		if ts == nil {
			c.AbortWithStatusJSON(401, utils.H{
				"code":    0,
				"message": "authentication has expired, please log in again",
			})
			return
		}

		//claims, err := x.IndexService.Verify(ctx, string(ts))
		//if err != nil {
		//	c.AbortWithStatusJSON(401, utils.H{
		//		"code":    0,
		//		"message": common.ErrAuthenticationExpired.Error(),
		//	})
		//	return
		//}
		//
		//c.Set("identity", claims)
		c.Next(ctx)
	}
}

func (x *API) Initialize(ctx context.Context) (h *server.Hertz, err error) {
	h = x.Hertz

	return
}
