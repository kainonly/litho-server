package app

import (
	"github.com/gin-gonic/gin"
	wpx "github.com/weplanx/go"
	"github.com/weplanx/go/helper"
	"github.com/weplanx/go/passport"
)

func authSystem(auth *passport.Auth, cookie *helper.CookieHelper) gin.HandlerFunc {
	return wpx.Returns(func(c *gin.Context) interface{} {
		tokenString, err := cookie.Get(c, "system_access_token")
		if err != nil {
			c.Abort()
			return err
		}
		claims, err := auth.Verify(tokenString)
		if err != nil {
			c.Abort()
			return err
		}
		c.Set("access_token", claims)
		c.Next()
		return nil
	})
}
