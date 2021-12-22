package roles

import (
	"api/common"
	"api/model"
	"github.com/gofiber/fiber/v2"
	"github.com/weplanx/go/api"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
)

type Controller struct {
	*InjectController
	*api.Controller
}

type InjectController struct {
	common.Inject
	Service *Service
}

func (x *Controller) Create(c *fiber.Ctx) interface{} {
	var body struct {
		Key    string              `bson:"key" json:"key"`
		Parent *primitive.ObjectID `bson:"parent" json:"parent"`
		Name   string              `bson:"name" json:"name"`
		Status *bool               `bson:"status" json:"status"`
	}
	if err := c.BodyParser(&body); err != nil {
		return err
	}
	data := model.NewRole(body.Key, body.Name)
	if body.Parent != nil {
		data.SetParent(data.Parent)
	}
	log.Println(body)
	result, err := x.API.Create(c.UserContext(), data)
	if err != nil {
		return err
	}
	return result
}
