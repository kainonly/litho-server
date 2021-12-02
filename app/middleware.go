package app

//func authSystem(auth *passport.Auth, cookie *helper.CookieHelper) fiber.Handler {
//	return wpx.Returns(func(c *fiber.Ctx) interface{} {
//		tokenString, err := cookie.Get(c, "system_access_token")
//		if err != nil {
//			c.Abort()
//			return err
//		}
//		claims, err := auth.Verify(tokenString)
//		if err != nil {
//			c.Abort()
//			return err
//		}
//		c.Set("access_token", claims)
//		c.Next()
//		return nil
//	})
//}
