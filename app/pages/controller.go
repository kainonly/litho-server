package pages

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Controller struct {
	Service *Service
}

func (x *Controller) SchemaKeyExists(c *gin.Context) interface{} {
	var query struct {
		Key string `form:"key" binding:"required"`
	}
	if err := c.ShouldBindQuery(&query); err != nil {
		return err
	}
	ctx := c.Request.Context()
	code, err := x.Service.SchemaKeyExists(ctx, query.Key)
	if err != nil {
		return err
	}
	return gin.H{
		"status": code,
	}
}

type ReorganizationDto struct {
	Id     primitive.ObjectID   `json:"id" binding:"required"`
	Parent primitive.ObjectID   `json:"parent" binding:"required"`
	Sort   []primitive.ObjectID `json:"sort" binding:"required"`
}

func (x *Controller) Reorganization(c *gin.Context) interface{} {
	var body ReorganizationDto
	if err := c.ShouldBindJSON(&body); err != nil {
		return err
	}
	ctx := c.Request.Context()
	result, err := x.Service.Reorganization(ctx, body)
	if err != nil {
		return err
	}
	return result
}

type SortSchemaFieldsDto struct {
	Id     primitive.ObjectID `json:"id" binding:"required"`
	Fields []string           `json:"fields" binding:"required"`
}

func (x *Controller) SortSchemaFields(c *gin.Context) interface{} {
	var body SortSchemaFieldsDto
	if err := c.ShouldBindJSON(&body); err != nil {
		return err
	}
	ctx := c.Request.Context()
	result, err := x.Service.SortSchemaFields(ctx, body)
	if err != nil {
		return err
	}
	return result
}

type DeleteSchemaFieldDto struct {
	Id  primitive.ObjectID `json:"id" binding:"required"`
	Key string             `json:"key" binding:"required"`
}

func (x *Controller) DeleteSchemaField(c *gin.Context) interface{} {
	var body DeleteSchemaFieldDto
	if err := c.ShouldBindJSON(&body); err != nil {
		return err
	}
	ctx := c.Request.Context()
	result, err := x.Service.DeleteSchemaField(ctx, body)
	if err != nil {
		return err
	}
	return result
}

type FindIndexesDto struct {
	Id primitive.ObjectID `json:"id" binding:"required"`
}

func (x *Controller) FindIndexes(c *gin.Context) interface{} {
	var body FindIndexesDto
	if err := c.ShouldBindJSON(&body); err != nil {
		return err
	}
	ctx := c.Request.Context()
	page, err := x.Service.FindOnePage(ctx, body.Id)
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
	Id     primitive.ObjectID `json:"id" binding:"required"`
	Name   string             `json:"name" binding:"required"`
	Keys   bson.D             `json:"keys" binding:"required,gt=0"`
	Unique *bool              `json:"unique" binding:"required"`
}

func (x *Controller) CreateIndex(c *gin.Context) interface{} {
	var body CreateIndexDto
	if err := c.ShouldBindJSON(&body); err != nil {
		return err
	}
	ctx := c.Request.Context()
	page, err := x.Service.FindOnePage(ctx, body.Id)
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
	Id   primitive.ObjectID `json:"id" binding:"required"`
	Name string             `json:"name" binding:"required"`
}

func (x *Controller) DeleteIndex(c *gin.Context) interface{} {
	var body DeleteIndexDto
	if err := c.ShouldBindJSON(&body); err != nil {
		return err
	}
	ctx := c.Request.Context()
	page, err := x.Service.FindOnePage(ctx, body.Id)
	if err != nil {
		return err
	}
	if _, err = x.Service.DeleteIndex(ctx, body, page.Schema.Key); err != nil {
		return err
	}
	return "ok"
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
