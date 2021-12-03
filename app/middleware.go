package app

import (
	"github.com/gofiber/fiber/v2"
	"github.com/weplanx/go/passport"
)

func AuthGuard(passport *passport.Passport) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenString := c.Cookies("console")
		claims, err := passport.Verify(tokenString)
		if err != nil {
			return err
		}
		c.Locals("claims", claims)
		return c.Next()
	}
}
