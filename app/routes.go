package app

import "github.com/gofiber/fiber/v2"

func Home(c *fiber.Ctx) interface{} {
	return fiber.Map{"msg": "hi"}
}
