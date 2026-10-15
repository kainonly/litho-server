package api

import (
	"context"
	"litho-api/common"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
)

func (x *API) Auth() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		ts := c.Cookie("TOKEN")
		if ts == nil {
			c.AbortWithStatusJSON(401, utils.H{
				"code":    0,
				"message": `认证已失效，请重新登录`,
			})
			return
		}

		claims, err := x.IndexX.Verify(ctx, string(ts))
		if err != nil {
			x.IndexX.ClearAccessToken(c)
			c.AbortWithStatusJSON(401, utils.H{
				"code":    0,
				"message": `认证已失效，请重新登录`,
			})
			return
		}

		// 获取登录用户信息
		var user *common.IAMUser
		if user, err = x.UsersX.GetIAMUser(ctx, claims.ActiveId); err != nil {
			x.IndexX.ClearAccessToken(c)
			c.AbortWithStatusJSON(401, utils.H{
				"code":    1,
				"message": `认证已失效，请重新登录`,
			})
			return
		}

		// 检测用户状态
		if !user.Status {
			x.IndexX.ClearAccessToken(c)
			c.AbortWithStatusJSON(401, utils.H{
				"code":    0,
				"message": `认证已失效，您的账户已被管理员禁用`,
			})
			return
		}

		user.Ip = string(c.GetHeader(x.V.Ip))

		// 未设置权限组
		if user.RoleID == "0" {
			x.IndexX.ClearAccessToken(c)
			c.AbortWithStatusJSON(403, utils.H{
				"code":    0,
				"message": `登录禁用，您的账户尚未设置权限组`,
			})
			return
		}

		// 获取用户唯一权限组的缓存策略
		if user.Strategy, err = x.RolesX.GetIAMRole(ctx, user.RoleID); err != nil {
			x.IndexX.ClearAccessToken(c)
			c.AbortWithStatusJSON(403, utils.H{
				"code":    0,
				"message": `登录禁用，您的账户尚未设置权限组`,
			})
			return
		}

		segments := strings.Split(string(c.GetRequest().Path())[1:], "/")
		c.Set(`segments`, segments)
		// 存在固定无需验证的接口
		if x.IgnoreRoutes(segments) {
			c.Set("identity", user)
			c.Next(ctx)
			return
		}

		// 获取用户访问路由对应缓存策略
		//u, _ := url.Parse(string(c.GetHeader(`x-route`)))
		//var route *common.IAMRoute
		//if route, err = x.RoutesX.GetIAMRoute(ctx, u.Path); err != nil {
		//	c.AbortWithStatusJSON(403, utils.H{
		//		"code":    0,
		//		"message": `禁止访问，您的账户不具备页面访问权限`,
		//	})
		//	return
		//}
		//
		//if !x.AllowRoute(route.ID, user.Strategy.Routes) {
		//	c.AbortWithStatusJSON(403, utils.H{
		//		"code":    0,
		//		"message": `禁止访问，您的账户不具备页面访问权限`,
		//	})
		//	return
		//}

		// 验证路由接口权限（已实现，开发时忽略）
		//if !x.VerifyACL(string(c.Method()), segments, route.Strategy) {
		//	c.AbortWithStatusJSON(403, utils.H{
		//		"code":    0,
		//		"message": `禁止访问，您的账户不具备页面接口访问权限`,
		//	})
		//	return
		//}

		//fmt.Println(u) // 访问页面

		c.Set("identity", user)
		c.Next(ctx)
	}
}

var Ignore = map[string]bool{
	`user`:   true,
	`layout`: true,
}

func (x *API) IgnoreRoutes(segments []string) bool {
	return Ignore[segments[0]]
}

func (x *API) AllowRoute(id string, routes common.A) bool {
	for _, route := range routes {
		if route.(string) == id {
			return true
		}
	}
	return false
}

func (x *API) VerifyACL(method string, segments []string, acl common.M) bool {
	if method == "GET" {
		return true
	}
	if method != "POST" || len(segments) < 2 {
		return false
	}
	key := segments[0]
	var strategy []any
	switch values := acl[key].(type) {
	case []any:
		strategy = values
	case common.A:
		strategy = values
	default:
		return false
	}
	strategyM := make(map[string]bool)
	for _, v := range strategy {
		action, ok := v.(string)
		if !ok {
			return false
		}
		// 存在允许完全访问的接口
		if action == "*" {
			return true
		}
		strategyM[action] = true
	}
	act := segments[1]
	if strategyM[act] {
		return true
	}
	return false
}
