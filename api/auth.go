package api

import (
	"context"
	"server/common"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
)

func (x *API) Auth() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		ts := c.Cookie("ACCESS_TOKEN")
		if ts == nil {
			c.AbortWithStatusJSON(401, common.M{
				"code":    0,
				"message": `身份验证已过期，请重新登录`,
			})
			return
		}

		claims, err := x.IndexX.Verify(ctx, string(ts))
		if err != nil {
			x.IndexX.ClearAccessToken(c)
			c.AbortWithStatusJSON(401, common.M{
				"code":    0,
				"message": `身份验证已过期，请重新登录`,
			})
			return
		}

		// 获取登录用户信息
		var user *common.IAMUser
		if user, err = x.UsersX.GetIAMUser(ctx, claims.ActiveId); err != nil {
			x.IndexX.ClearAccessToken(c)
			c.AbortWithStatusJSON(401, common.M{
				"code":    0,
				"message": `身份验证已过期，请重新登录`,
			})
			return
		}

		// 检测企业成员状态
		if !user.Active {
			x.IndexX.ClearAccessToken(c)
			c.AbortWithStatusJSON(401, common.M{
				"code":    0,
				"message": `身份验证已过期，您的账号已被管理员禁用`,
			})
			return
		}

		// 未设置权限组
		if user.RoleID == "0" {
			x.IndexX.ClearAccessToken(c)
			c.AbortWithStatusJSON(403, common.M{
				"code":    0,
				"message": `禁止登录，您的账号未分配角色`,
			})
			return
		}

		// 获取用户对应权限缓存策略
		if user.Strategy, err = x.RolesX.GetIAMRole(ctx, user.RoleID); err != nil {
			x.IndexX.ClearAccessToken(c)
			c.AbortWithStatusJSON(403, utils.H{
				"code":    0,
				"message": `登录禁用，您的账户尚未设置权限组`,
			})
			return
		}

		c.Set("identity", user)
		c.Next(ctx)
	}
}
