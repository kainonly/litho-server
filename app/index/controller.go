package index

import (
	"api/common"
	"github.com/gofiber/fiber/v2"
)

type Controller struct {
	*InjectController
}

type InjectController struct {
	common.Inject
	Service *Service
}

func NewController(i *InjectController) *Controller {
	return &Controller{
		InjectController: i,
	}
}

func (x *Controller) Index(c *fiber.Ctx) interface{} {
	c.Cookie(&fiber.Cookie{
		Name:     "name",
		Value:    "kain",
		Secure:   true,
		HTTPOnly: true,
		SameSite: "Strict",
	})
	return fiber.Map{
		"msg": "hi",
	}
}

func (x *Controller) Get(c *fiber.Ctx) interface{} {
	return fiber.Map{
		"name": c.Cookies("name"),
	}
}
