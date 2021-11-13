package app

import (
	"api/common"
	"github.com/gin-gonic/gin"
	"github.com/weplanx/go/helper"
	"github.com/weplanx/go/passport"
)

func authSystem(auth *passport.Auth, cookie *helper.CookieHelper) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := cookie.Get(c, "system_access_token")
		if err != nil {
			c.AbortWithStatusJSON(200, gin.H{
				"code":    1,
				"message": err.Error(),
			})
			return
		}
		claims, err := auth.Verify(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(200, gin.H{
				"code":    1,
				"message": common.LoginExpired,
			})
			return
		}
		c.Set("access_token", claims)
		c.Next()
	}
}
