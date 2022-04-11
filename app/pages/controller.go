package pages

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Controller struct {
	Service *Service
}

func (x *Controller) Sort(c *gin.Context) interface{} {
	var body struct {
		Sort []primitive.ObjectID `json:"sort" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		return err
	}
	ctx := c.Request.Context()
	result, err := x.Service.Sort(ctx, body.Sort)
	if err != nil {
		return err
	}
	return result
}

func (x *Controller) Indexes(c *gin.Context) interface{} {
	var params struct {
		Id string `uri:"id" binding:"required,objectId"`
	}
	if err := c.ShouldBindUri(&params); err != nil {
		return err
	}
	ctx := c.Request.Context()
	data, err := x.Service.FindOneById(ctx, params.Id)
	if err != nil {
		return err
	}
	indexes, err := x.Service.FindIndexes(ctx, data.Schema.Key)
	if err != nil {
		return err
	}
	return indexes
}

func (x *Controller) CreateIndex(c *gin.Context) interface{} {
	var params struct {
		Id   string `uri:"id" binding:"required,objectId"`
		Name string `uri:"name" binding:"required,key"`
	}
	if err := c.ShouldBindUri(&params); err != nil {
		return err
	}
	var body struct {
		Keys   bson.D `json:"keys" binding:"required,gt=0"`
		Unique *bool  `json:"unique" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		return err
	}
	ctx := c.Request.Context()
	page, err := x.Service.FindOneById(ctx, params.Id)
	if err != nil {
		return err
	}
	if _, err = x.Service.CreateIndex(ctx, page.Schema.Key, params.Name, body.Keys, *body.Unique); err != nil {
		return err
	}
	return nil
}

func (x *Controller) DeleteIndex(c *gin.Context) interface{} {
	var params struct {
		Id   string `uri:"id" binding:"required,objectId"`
		Name string `uri:"name" binding:"required,key"`
	}
	if err := c.ShouldBindUri(&params); err != nil {
		return err
	}
	ctx := c.Request.Context()
	page, err := x.Service.FindOneById(ctx, params.Id)
	if err != nil {
		return err
	}
	if _, err = x.Service.DeleteIndex(ctx, page.Schema.Key, params.Name); err != nil {
		return err
	}
	return nil
}
