package api

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/hertz-contrib/jwt"
)

func (x *API) Routes(h *server.Hertz) (auth *jwt.HertzJWTMiddleware, err error) {
	//if auth, err = x.Auth(); err != nil {
	//	return
	//}

	h.GET("", x.IndexController.Index)
	h.POST("login", auth.LoginHandler)
	h.GET("code", auth.MiddlewareFunc(), x.IndexController.GetRefreshCode)
	h.POST("refresh_token", auth.MiddlewareFunc(), x.IndexController.VerifyRefreshCode, auth.RefreshHandler)
	h.POST("logout", auth.MiddlewareFunc(), auth.LogoutHandler)

	h.GET("navs", auth.MiddlewareFunc(), x.IndexController.GetNavs)
	h.GET("options", auth.MiddlewareFunc(), x.IndexController.GetOptions)

	_user := h.Group("user", auth.MiddlewareFunc())
	{
		_user.GET("", x.IndexController.GetUser)
		_user.PATCH("", x.IndexController.SetUser)
	}

	//_values := h.Group("values")
	//{
	//	_values.GET("", x.ValuesController.Get)
	//	_values.PATCH("", x.ValuesController.Set)
	//	_values.DELETE(":key", x.ValuesController.Remove)
	//}
	//
	//_sessions := h.Group("sessions", auth.MiddlewareFunc())
	//{
	//	_sessions.GET("", x.SessionController.Lists)
	//	_sessions.DELETE(":uid", x.SessionController.Remove)
	//	_sessions.DELETE("", x.SessionController.Clear)
	//}
	//
	//_dsl := h.Group("/:model", auth.MiddlewareFunc())
	//{
	//	_dsl.POST("", x.DslController.Create)
	//	_dsl.POST("bulk-create", x.DslController.BulkCreate)
	//	_dsl.GET("_size", x.DslController.Size)
	//	_dsl.GET("", x.DslController.Find)
	//	_dsl.GET("_one", x.DslController.FindOne)
	//	_dsl.GET(":id", x.DslController.FindById)
	//	_dsl.PATCH("", x.DslController.Update)
	//	_dsl.PATCH(":id", x.DslController.UpdateById)
	//	_dsl.PUT(":id", x.DslController.Replace)
	//	_dsl.DELETE(":id", x.DslController.Delete)
	//	_dsl.POST("bulk-delete", x.DslController.BulkDelete)
	//	_dsl.POST("sort", x.DslController.Sort)
	//}

	//_pages := h.Group("pages", auth.MiddlewareFunc())
	//{
	//}

	return
}
