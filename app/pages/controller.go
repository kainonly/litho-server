package pages

import (
	"api/common"
	"github.com/gin-gonic/gin"
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
	Service *Service
}

func (x *Controller) Delete(c *gin.Context) interface{} {
	result := x.Controller.Delete(c)
	if _, ok := result.(error); ok {
		return result
	}
	// TODO: 发送变更集合名队列
	return result
}

type CheckKeyDto struct {
	Value string `json:"value" binding:"required"`
}

func (x *Controller) CheckKey(c *gin.Context) interface{} {
	var body CheckKeyDto
	if err := c.ShouldBindJSON(&body); err != nil {
		return err
	}
	ctx := c.Request.Context()
	result, err := x.Service.CheckKey(ctx, body)
	if err != nil {
		return err
	}
	return result
}

type ReorganizationDto struct {
	Id     primitive.ObjectID   `json:"id" binding:"required"`
	Parent primitive.ObjectID   `json:"parent" binding:"required"`
	Sort   []primitive.ObjectID `json:"sort" binding:"required"`
}

// Reorganization 重组
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
	Id   primitive.ObjectID `json:"id" binding:"required"`
	Name string             `json:"name" binding:"required"`
}

func (x *Controller) DeleteIndex(c *gin.Context) interface{} {
	var body DeleteIndexDto
	if err := c.ShouldBindJSON(&body); err != nil {
		return err
	}
	ctx := c.Request.Context()
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
