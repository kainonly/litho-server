package pages

import (
	"api/common"
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/thoas/go-funk"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Controller struct {
	*InjectController
}

type InjectController struct {
	common.Inject
	Service *Service
}

func (x *Controller) Sort(c *fiber.Ctx) interface{} {
	return nil
}

func (x *Controller) CheckKey(c *fiber.Ctx) interface{} {
	var body struct {
		Value string `json:"value" validate:"required"`
	}
	if err := c.BodyParser(&body); err != nil {
		return err
	}
	if err := validator.New().Struct(body); err != nil {
		return err
	}
	ctx := c.UserContext()
	count, err := x.Db.Collection("pages").CountDocuments(ctx, bson.M{
		"schema.key": body.Value,
	})
	if err != nil {
		return err
	}
	if count != 0 {
		return "duplicated"
	}
	collections, err := x.Db.ListCollectionNames(ctx, bson.M{})
	if err != nil {
		return err
	}
	if funk.Contains(collections, body.Value) {
		return "history"
	}
	return "ok"
}

func (x *Controller) SortFields(c *fiber.Ctx) interface{} {
	var body struct {
		Id     primitive.ObjectID `json:"id" validate:"required"`
		Fields bson.A             `json:"fields" validate:"required"`
	}
	if err := c.BodyParser(&body); err != nil {
		return err
	}
	if err := validator.New().Struct(body); err != nil {
		return err
	}
	result, err := x.Db.Collection("pages").UpdateOne(context.TODO(), bson.M{
		"_id": body.Id,
	}, bson.M{
		"$set": bson.M{
			"schema.fields": body.Fields,
		},
	})
	if err != nil {
		return err
	}
	return result
}
