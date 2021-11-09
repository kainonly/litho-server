package schema

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/thoas/go-funk"
	"github.com/weplanx/support/api"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"laboratory/common"
)

type InjectController struct {
	*common.App
	Service *Service
}

type Controller struct {
	*InjectController
	*api.API
}

func NewController(i *InjectController) *Controller {
	return &Controller{
		InjectController: i,
		API: api.New(
			i.Mongo,
			i.Db,
			api.SetCollection("schema"),
			api.ProjectionNone(),
		),
	}
}

func (x *Controller) ExistsCollection(c *gin.Context) interface{} {
	var body struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		return err
	}
	count, err := x.Db.Collection("schema").CountDocuments(c, bson.M{
		"collection": body.Name,
	})
	if err != nil {
		return err
	}
	if count != 0 {
		return true
	}
	collections, err := x.Db.ListCollectionNames(c, bson.M{})
	if err != nil {
		return err
	}
	return funk.Contains(collections, body.Name)
}

func (x *Controller) isSystem(c *gin.Context, id string) error {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	var data map[string]interface{}
	if err := x.Db.Collection("schema").FindOne(c, bson.M{
		"_id": objectId,
	}).Decode(&data); err != nil {
		return err
	}
	if data["system"] == true {
		return errors.New("该集合为系统组件，不可删除修改")
	}
	return nil
}

func (x *Controller) Update(c *gin.Context) interface{} {
	var body struct {
		api.UpdateBody
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		return err
	}
	h := api.StartHook(c)
	h.SetBody(body)
	if err := x.isSystem(c, body.Where["_id"].(string)); err != nil {
		return err
	}
	return x.API.Update(c)
}

func (x *Controller) Delete(c *gin.Context) interface{} {
	var body struct {
		api.DeleteBody
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		return err
	}
	h := api.StartHook(c)
	h.SetBody(body)
	if err := x.isSystem(c, body.Where["_id"].(string)); err != nil {
		return err
	}
	return x.API.Delete(c)
}

func (x *Controller) Sort(c *gin.Context) interface{} {
	var body struct {
		Id     primitive.ObjectID `json:"id" binding:"required"`
		Fields bson.A             `json:"fields" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		return err
	}
	result, err := x.Db.Collection("schema").UpdateOne(c, bson.M{
		"_id": body.Id,
	}, bson.M{
		"$set": bson.M{
			"fields": body.Fields,
		},
	})
	if err != nil {
		return err
	}
	return result
}
