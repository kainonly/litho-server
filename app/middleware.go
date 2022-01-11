package app

import (
	"api/common"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/requestid"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/weplanx/go/engine"
	"github.com/weplanx/go/passport"
	"go.uber.org/zap"
	"time"
)

func globalMiddleware(r *gin.Engine, values *common.Values) *gin.Engine {
	r.SetTrustedProxies(values.TrustedProxies)
	logger, _ := zap.NewProduction()
	r.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	r.Use(gin.CustomRecovery(catchError))
	r.Use(requestid.New())
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

func catchError(c *gin.Context, err interface{}) {
	c.AbortWithStatusJSON(500, gin.H{
		"message": err,
	})
}

func authGuard(passport *passport.Passport) gin.HandlerFunc {
	return func(c *gin.Context) {
		ts, err := c.Cookie("access_token")
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{
				"code":    "AUTH_EXPIRED",
				"message": common.LoginExpired.Error(),
			})
			return
		}
		claims, err := passport.Verify(ts)
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
