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

type Controller struct {
	*InjectController
	*api.API
}

type InjectController struct {
	common.App
	Service *Service
}

func (x *Controller) ExistsCollection(c *gin.Context) interface{} {
	var body struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		return err
	}
	count, err := x.Collection.CountDocuments(c, bson.M{
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

func (x *Controller) Delete(c *gin.Context) interface{} {
	var body struct {
		api.DeleteBody
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		return err
	}
	h := api.StartHook(c)
	h.SetBody(body)
	objectId, err := primitive.ObjectIDFromHex(body.Where["_id"].(string))
	if err != nil {
		return err
	}
	var data map[string]interface{}
	if err := x.Collection.FindOne(c, bson.M{
		"_id": objectId,
	}).Decode(&data); err != nil {
		return err
	}
	if data["system"] == true {
		return errors.New("该集合为系统组件，不可删除修改")
	}
	return x.API.Delete(c)
}
