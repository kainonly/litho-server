package devops

import (
	"api/common"
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/weplanx/go/basic"
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

func (x *Controller) Setup(c *fiber.Ctx) interface{} {
	ctx := context.Background()
	if err := basic.GenerateSchema(ctx, x.Db); err != nil {
		return err
	}
	if err := basic.GeneratePage(ctx, x.Db); err != nil {
		return err
	}
	if err := basic.GenerateRoleAndAdmin(ctx, x.Db); err != nil {
		return err
	}
	return fiber.Map{"message": "ok"}
}
