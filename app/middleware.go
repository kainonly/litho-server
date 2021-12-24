package app

import (
	"api/common"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/weplanx/go/api"
	"time"
)

func middleware(r *gin.Engine, values *common.Values) *gin.Engine {
	api.RegisterValidation()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     values.Cors.AllowOrigins,
		AllowMethods:     values.Cors.AllowMethods,
		AllowHeaders:     values.Cors.AllowHeaders,
		AllowCredentials: values.Cors.AllowCredentials,
		MaxAge:           time.Duration(values.Cors.MaxAge) * time.Second,
	}))
	return r
}

//func AuthGuard(passport *passport.Passport) fiber.Handler {
//	return func(c *fiber.Ctx) error {
//		tokenString := c.Cookies("access_token")
//		if tokenString == "" {
//			return c.JSON(fiber.Map{
//				"code":    401,
//				"message": common.LoginExpired.Error(),
//			})
//		}
//		claims, err := passport.Verify(tokenString)
//		if err != nil {
//			return err
//		}
//		c.Locals(common.TokenClaimsKey, claims)
//		return c.Next()
//	}
//}
