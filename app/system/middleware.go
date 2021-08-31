package system

import (
	"github.com/gin-gonic/gin"
	"github.com/kainonly/go-bit/authx"
	"github.com/kainonly/go-bit/cookie"
)

func authMiddleware(auth *authx.Auth, cookie *cookie.Cookie) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := cookie.Get(c, "system_access_token")
		if err != nil {
			c.AbortWithStatusJSON(200, gin.H{
				"error": 0,
				"msg":   err.Error(),
			})
			return
		}
		claims, err := auth.Verify(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(200, gin.H{
				"error": 0,
				"msg":   err.Error(),
			})
			return
		}
		c.Set("access_token", claims)
		c.Next()
	}
}
