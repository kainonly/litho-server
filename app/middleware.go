package app

import (
	"api/common"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/weplanx/go/engine"
	"github.com/weplanx/go/passport"
	"time"
)

func globalMiddleware(r *gin.Engine, values *common.Values) *gin.Engine {
	r.SetTrustedProxies(values.TrustedProxies)
	r.Use(cors.New(cors.Config{
		AllowOrigins:     values.Cors.AllowOrigins,
		AllowMethods:     values.Cors.AllowMethods,
		AllowHeaders:     values.Cors.AllowHeaders,
		ExposeHeaders:    values.Cors.ExposeHeaders,
		AllowCredentials: values.Cors.AllowCredentials,
		MaxAge:           time.Duration(values.Cors.MaxAge) * time.Second,
	}))
	engine.RegisterValidation()
	return r
}

func authGuard(passport *passport.Passport) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("access_token")
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{
				"code":    "AUTH_EXPIRED",
				"message": common.LoginExpired.Error(),
			})
			return
		}
		claims, err := passport.Verify(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{
				"code":    "AUTH_EXPIRED",
				"message": common.LoginExpired.Error(),
			})
			return
		}
		c.Set(common.TokenClaimsKey, claims)
		c.Next()
	}
}
