package app

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
