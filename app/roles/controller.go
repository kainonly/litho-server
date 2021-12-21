package roles

import (
	"api/common"
	"api/model"
	"github.com/gofiber/fiber/v2"
	"github.com/weplanx/go/api"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Controller struct {
	*InjectController
	*api.Controller
}

type InjectController struct {
	common.Inject
	APIs    *api.API
	Service *Service
}

func (x *Controller) Create(c *fiber.Ctx) interface{} {
	var body model.Role
	if err := c.BodyParser(&body); err != nil {
		return err
	}
	var parent primitive.ObjectID
	if body.Parent != "root" {
		var err error
		if parent, err = primitive.ObjectIDFromHex(body.Parent.(string)); err != nil {
			return err
		}
	}
	result, err := x.APIs.Create(
		c.UserContext(),
		model.NewRole(body.Key, body.Name).SetParent(parent),
	)
	if err != nil {
		return err
	}
	return result
}
