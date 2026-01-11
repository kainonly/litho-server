package api

import (
	"context"
	"server/common"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
)

func (x *API) Auth() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		ts := c.Cookie("TOKEN")
		if ts == nil {
			c.AbortWithStatusJSON(401, utils.H{
				"code":    0,
				"message": `Authentication expired, please login again`,
			})
			return
		}

		claims, err := x.IndexX.Verify(ctx, string(ts))
		if err != nil {
			x.IndexX.ClearAccessToken(c)
			c.AbortWithStatusJSON(401, utils.H{
				"code":    0,
				"message": `Authentication expired, please login again`,
			})
			return
		}

		// 获取登录用户信息
		var user *common.IAMUser
		if user, err = x.UsersX.GetIAMUser(ctx, claims.ActiveId); err != nil {
			x.IndexX.ClearAccessToken(c)
			c.AbortWithStatusJSON(401, utils.H{
				"code":    0,
				"message": `Authentication expired, please login again`,
			})
			return
		}

		// 检测企业成员状态
		if !user.Active {
			x.IndexX.ClearAccessToken(c)
			c.AbortWithStatusJSON(401, utils.H{
				"code":    0,
				"message": `Authentication expired, your account has been disabled by administrator`,
			})
			return
		}

		// 未设置权限组
		if user.RoleID == "0" {
			x.IndexX.ClearAccessToken(c)
			c.AbortWithStatusJSON(403, utils.H{
				"code":    0,
				"message": `Login disabled, your account has no role assigned`,
			})
			return
		}

		c.Set("identity", user)
		c.Next(ctx)
	}
}
