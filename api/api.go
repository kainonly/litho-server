package api

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/google/wire"
	"github.com/weplanx/go/csrf"
	"github.com/weplanx/server/api/index"
	"github.com/weplanx/server/common"
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

func (x *API) Routes(h *server.Hertz) (err error) {
	//csrfToken := x.Csrf.VerifyToken(!x.V.IsRelease())
	//auth := x.AuthGuard()
	//audit := x.Audit()

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

		//claims, err := x.IndexX.Verify(ctx, string(ts))
		//if err != nil {
		//	common.ClearAccessToken(c)
		//	c.AbortWithStatusJSON(401, utils.H{
		//		"code":    0,
		//		"message": common.ErrAuthenticationExpired.Error(),
		//	})
		//	return
		//}

		//c.Set("identity", claims)
		c.Next(ctx)
	}
}

//func (x *API) Audit() app.HandlerFunc {
//	return func(ctx context.Context, c *app.RequestContext) {
//		now := time.Now()
//		c.Next(ctx)
//		method := string(c.Request.Header.Method())
//		if method == "GET" {
//			return
//		}
//		var userId string
//		if value, ok := c.Get("identity"); ok {
//			claims := value.(passport.Claims)
//			userId = claims.UserId
//		}
//
//		format := map[string]interface{}{
//			"body": "json",
//		}
//		if userId != "" {
//			format["metadata.user_id"] = "oid"
//		}
//	}
//}

func (x *API) Initialize(ctx context.Context) (h *server.Hertz, err error) {
	h = x.Hertz

	return
}
