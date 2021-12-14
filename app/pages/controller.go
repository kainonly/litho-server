package pages

import (
	"api/common"
	"context"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/thoas/go-funk"
	"github.com/weplanx/go/api"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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

// Reorganization 重组
func (x *Controller) Reorganization(c *fiber.Ctx) interface{} {
	var body struct {
		Id     *primitive.ObjectID   `json:"id" validate:"required"`
		Parent string                `json:"parent" validate:"required"`
		Sort   []*primitive.ObjectID `json:"sort" validate:"required"`
	}
	if err := c.BodyParser(&body); err != nil {
		return err
	}
	ctx := c.UserContext()
	models := []mongo.WriteModel{
		mongo.NewUpdateOneModel().
			SetFilter(bson.M{"_id": body.Id}).
			SetUpdate(bson.M{"$set": bson.M{"parent": body.Parent}}),
	}
	for i, x := range body.Sort {
		models = append(models, mongo.NewUpdateOneModel().
			SetFilter(bson.M{"_id": x}).
			SetUpdate(bson.M{"$set": bson.M{"sort": i}}),
		)
	}
	result, err := x.Db.Collection("pages").BulkWrite(ctx, models)
	if err != nil {
		return err
	}
	return result
}

func (x *Controller) SortSchemaFields(c *fiber.Ctx) interface{} {
	var body struct {
		Id     primitive.ObjectID `json:"id" validate:"required"`
		Fields []string           `json:"fields" validate:"required"`
	}
	if err := c.BodyParser(&body); err != nil {
		return err
	}
	if err := validator.New().Struct(body); err != nil {
		return err
	}
	values := make(bson.M, 0)
	for i, x := range body.Fields {
		values[fmt.Sprintf(`schema.fields.%s.sort`, x)] = i
	}
	result, err := x.Db.Collection("pages").
		UpdateOne(context.TODO(),
			bson.M{"_id": body.Id},
			bson.M{"$set": values},
		)
	if err != nil {
		return err
	}
	return result
}

func (x *Controller) DeleteSchemaField(c *fiber.Ctx) interface{} {
	var body struct {
		Id  primitive.ObjectID `json:"id" validate:"required"`
		Key string             `json:"key" validate:"required"`
	}
	if err := c.BodyParser(&body); err != nil {
		return err
	}
	if err := validator.New().Struct(body); err != nil {
		return err
	}
	result, err := x.Db.Collection("pages").
		UpdateOne(context.TODO(),
			bson.M{"_id": body.Id},
			bson.M{"$unset": bson.M{
				fmt.Sprintf("schema.fields.%s", body.Key): "",
			}},
		)
	if err != nil {
		return err
	}
	return result
}
