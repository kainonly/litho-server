package pages

import (
	"api/common"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/weplanx/go/api"
	"go.mongodb.org/mongo-driver/bson"
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

func (x *Controller) Delete(c *fiber.Ctx) interface{} {
	result := x.Controller.Delete(c)
	if _, ok := result.(error); ok {
		return result
	}
	// TODO: 发送变更集合名队列
	return result
}

type CheckKeyDto struct {
	Value string `json:"value" validate:"required"`
}

func (x *Controller) CheckKey(c *fiber.Ctx) interface{} {
	var body CheckKeyDto
	if err := c.BodyParser(&body); err != nil {
		return err
	}
	if err := validator.New().Struct(body); err != nil {
		return err
	}
	ctx := c.UserContext()
	result, err := x.Service.CheckKey(ctx, body)
	if err != nil {
		return err
	}
	return result
}

type ReorganizationDto struct {
	Id     primitive.ObjectID   `json:"id" validate:"required"`
	Parent primitive.ObjectID   `json:"parent" validate:"required"`
	Sort   []primitive.ObjectID `json:"sort" validate:"required"`
}

// Reorganization 重组
func (x *Controller) Reorganization(c *fiber.Ctx) interface{} {
	var body ReorganizationDto
	if err := c.BodyParser(&body); err != nil {
		return err
	}
	ctx := c.UserContext()
	result, err := x.Service.Reorganization(ctx, body)
	if err != nil {
		return err
	}
	return result
}

type SortSchemaFieldsDto struct {
	Id     primitive.ObjectID `json:"id" validate:"required"`
	Fields []string           `json:"fields" validate:"required"`
}

func (x *Controller) SortSchemaFields(c *fiber.Ctx) interface{} {
	var body SortSchemaFieldsDto
	if err := c.BodyParser(&body); err != nil {
		return err
	}
	if err := validator.New().Struct(body); err != nil {
		return err
	}
	ctx := c.UserContext()
	result, err := x.Service.SortSchemaFields(ctx, body)
	if err != nil {
		return err
	}
	return result
}

type DeleteSchemaFieldDto struct {
	Id  primitive.ObjectID `json:"id" validate:"required"`
	Key string             `json:"key" validate:"required"`
}

func (x *Controller) DeleteSchemaField(c *fiber.Ctx) interface{} {
	var body DeleteSchemaFieldDto
	if err := c.BodyParser(&body); err != nil {
		return err
	}
	if err := validator.New().Struct(body); err != nil {
		return err
	}
	ctx := c.UserContext()
	result, err := x.Service.DeleteSchemaField(ctx, body)
	if err != nil {
		return err
	}
	return result
}

type FindIndexesDto struct {
	Id primitive.ObjectID `json:"id" validate:"required"`
}

func (x *Controller) FindIndexes(c *fiber.Ctx) interface{} {
	var body FindIndexesDto
	if err := c.BodyParser(&body); err != nil {
		return err
	}
	if err := validator.New().Struct(body); err != nil {
		return err
	}
	ctx := c.UserContext()
	page, err := x.Service.FindOneById(ctx, body.Id)
	if err != nil {
		return err
	}
	result, err := x.Service.FindIndexes(ctx, page.Schema.Key)
	if err != nil {
		return err
	}
	return result
}

type CreateIndexDto struct {
	Id     primitive.ObjectID `json:"id" validate:"required"`
	Name   string             `json:"name" validate:"required"`
	Keys   bson.D             `json:"keys" validate:"required,gt=0"`
	Unique *bool              `json:"unique" validate:"required"`
}

func (x *Controller) CreateIndex(c *fiber.Ctx) interface{} {
	var body CreateIndexDto
	if err := c.BodyParser(&body); err != nil {
		return err
	}
	if err := validator.New().Struct(body); err != nil {
		return err
	}
	ctx := c.UserContext()
	page, err := x.Service.FindOneById(ctx, body.Id)
	if err != nil {
		return err
	}
	result, err := x.Service.CreateIndex(ctx, body, page.Schema.Key)
	if err != nil {
		return err
	}
	return result
}

type DeleteIndexDto struct {
	Id   primitive.ObjectID `json:"id" validate:"required"`
	Name string             `json:"name" validate:"required"`
}

func (x *Controller) DeleteIndex(c *fiber.Ctx) interface{} {
	var body DeleteIndexDto
	if err := c.BodyParser(&body); err != nil {
		return err
	}
	if err := validator.New().Struct(body); err != nil {
		return err
	}
	ctx := c.UserContext()
	page, err := x.Service.FindOneById(ctx, body.Id)
	if err != nil {
		return err
	}
	if _, err = x.Service.DeleteIndex(ctx, body, page.Schema.Key); err != nil {
		return err
	}
	return "ok"
}

type UpdateValidatorDto struct {
	Id        primitive.ObjectID `json:"id" validate:"required"`
	Validator string             `json:"validator" validate:"required"`
}

func (x *Controller) UpdateValidator(c *fiber.Ctx) interface{} {
	var body UpdateValidatorDto
	if err := c.BodyParser(&body); err != nil {
		return err
	}
	if err := validator.New().Struct(body); err != nil {
		return err
	}
	ctx := c.UserContext()
	result, err := x.Service.UpdateValidator(ctx, body)
	if err != nil {
		return err
	}
	return result
}
