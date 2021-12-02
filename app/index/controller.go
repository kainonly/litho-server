package index

import (
	"api/common"
	"github.com/gofiber/fiber/v2"
)

type Controller struct {
	*InjectController
}

type InjectController struct {
	common.App
	Service *Service
}

func NewController(i *InjectController) *Controller {
	return &Controller{
		InjectController: i,
	}
}

func (x *Controller) Index(c *fiber.Ctx) interface{} {
	return fiber.Map{
		"msg": "hi",
	}
}
