package system

import (
	"api/common"
	"github.com/gin-gonic/gin"
	"github.com/weplanx/go/passport"
)

type Middleware struct {
	*Service
	Passport *passport.Passport
}

func (x *Middleware) AuthGuard() gin.HandlerFunc {
	return func(c *gin.Context) {
		ts, err := c.Cookie("access_token")
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{
				"code":    "AUTH_EXPIRED",
				"message": common.AuthExpired.Error(),
			})
			return
		}
		claims, err := x.Passport.Verify(ts)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{
				"code":    "AUTH_EXPIRED",
				"message": common.AuthExpired.Error(),
			})
			return
		}
		ctx := c.Request.Context()
		uid := claims["context"].(map[string]interface{})["uid"].(string)
		ok, err := x.VerifySession(ctx, uid, claims["jti"].(string))
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{
				"code":    "AUTH_CONFLICT",
				"message": common.AuthExpired.Error(),
			})
			return
		}
		if !ok {
			c.AbortWithStatusJSON(401, gin.H{
				"code":    "AUTH_CONFLICT",
				"message": common.AuthConflict.Error(),
			})
			return
		}
		if err = x.RenewSession(ctx, uid); err != nil {
			c.AbortWithStatusJSON(500, gin.H{
				"message": err.Error(),
			})
			return
		}
		c.Set(common.TokenClaimsKey, claims)
		c.Next()
	}
}
