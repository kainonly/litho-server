package pages

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Controller struct {
	Service *Service
}

func (x *Controller) HasSchemaKey(c *gin.Context) interface{} {
	var query struct {
		Key string `form:"key" binding:"required,key"`
	}
	if err := c.ShouldBindQuery(&query); err != nil {
		return err
	}
	ctx := c.Request.Context()
	code, err := x.Service.HasSchemaKey(ctx, query.Key)
	if err != nil {
		return err
	}
	return gin.H{
		"status": code,
	}
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

func (x *Controller) FindIndexes(c *gin.Context) interface{} {
	var params struct {
		Id string `uri:"id" binding:"required,objectId"`
	}
	if err := c.ShouldBindUri(&params); err != nil {
		return err
	}
	oid, _ := primitive.ObjectIDFromHex(params.Id)
	ctx := c.Request.Context()
	data, err := x.Service.FindOnePage(ctx, oid)
	if err != nil {
		return err
	}
	result, err := x.Service.FindIndexes(ctx, data.Schema.Key)
	if err != nil {
		return err
	}
	return result
}

type CreateIndexDto struct {
	Keys   bson.D `json:"keys" binding:"required,gt=0"`
	Unique *bool  `json:"unique" binding:"required"`
}

func (x *Controller) CreateIndex(c *gin.Context) interface{} {
	var params struct {
		Id   string `uri:"id" binding:"required,objectId"`
		Name string `uri:"name" binding:"required,key"`
	}
	if err := c.ShouldBindUri(&params); err != nil {
		return err
	}
	var body CreateIndexDto
	if err := c.ShouldBindJSON(&body); err != nil {
		return err
	}
	ctx := c.Request.Context()
	oid, _ := primitive.ObjectIDFromHex(params.Id)
	page, err := x.Service.FindOnePage(ctx, oid)
	if err != nil {
		return err
	}
	if _, err = x.Service.CreateIndex(ctx, page.Schema.Key, params.Name, body); err != nil {
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
	oid, _ := primitive.ObjectIDFromHex(params.Id)
	page, err := x.Service.FindOnePage(ctx, oid)
	if err != nil {
		return err
	}
	if _, err = x.Service.DeleteIndex(ctx, page.Schema.Key, params.Name); err != nil {
		return err
	}
	return nil
}

type UpdateValidatorDto struct {
	Id        primitive.ObjectID `json:"id" binding:"required"`
	Validator string             `json:"validator" binding:"required"`
}

func (x *Controller) UpdateValidator(c *gin.Context) interface{} {
	var body UpdateValidatorDto
	if err := c.ShouldBindJSON(&body); err != nil {
		return err
	}
	ctx := c.Request.Context()
	result, err := x.Service.UpdateValidator(ctx, body)
	if err != nil {
		return err
	}
	return result
}
