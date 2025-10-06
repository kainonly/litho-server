package api

import (
	"server/common"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"golang.org/x/net/context"
)

func (x *API) Auth() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		ts := c.Cookie("TOKEN")
		if ts == nil {
			c.AbortWithStatusJSON(401, utils.H{
				"code":    0,
				"message": `Authentication expired. Please re-login.`,
			})
			return
		}

		claims, err := x.IndexX.Verify(ctx, string(ts))
		if err != nil {
			x.IndexX.ClearAccessToken(c)
			c.AbortWithStatusJSON(401, utils.H{
				"code":    0,
				"message": `Authentication expired. Please re-login.`,
			})
			return
		}

		var user *common.IAMUser
		if user, err = x.UsersX.GetIAMUser(ctx, claims.ActiveId); err != nil {
			x.IndexX.ClearAccessToken(c)
			c.AbortWithStatusJSON(401, utils.H{
				"code":    0,
				"message": `Authentication expired. Please re-login.`,
			})
			return
		}

		if !*user.Status {
			x.IndexX.ClearAccessToken(c)
			c.AbortWithStatusJSON(401, utils.H{
				"code":    0,
				"message": `Authentication expired. Please re-login.`,
			})
			return
		}

		c.Set("identity", user)
		c.Next(ctx)
	}
}
