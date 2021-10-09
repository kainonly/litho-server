package schema

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/thoas/go-funk"
	"github.com/weplanx/support/api"
	"github.com/weplanx/support/basic"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"laboratory/common"
	"log"
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
	collections, err := x.Db.ListCollectionNames(c, bson.M{})
	if err != nil {
		return err
	}
	return funk.Contains(collections, body.Name)
}

func (x *Controller) Delete(c *gin.Context) interface{} {
	var body api.DeleteBody
	if err := c.ShouldBindJSON(&body); err != nil {
		return err
	}
	objectId, err := primitive.ObjectIDFromHex(body.Where["_id"].(string))
	if err != nil {
		return err
	}
	var data basic.Schema
	if err := x.Db.Collection("schema").FindOne(c, bson.M{
		"_id": objectId,
	}).Decode(&data); err != nil {
		return err
	}
	log.Println(data)
	if data.System == basic.True() {
		return errors.New("该集合为系统组件，不可删除修改")
	}
	return x.API.Delete(c)
}
