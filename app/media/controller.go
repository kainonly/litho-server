package media

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Controller struct {
	Service *Service
}

func (x *Controller) FindLabels(c *gin.Context) interface{} {
	ctx := c.Request.Context()
	values, err := x.Service.FindLabels(ctx)
	if err != nil {
		return err
	}
	return values
}

type BulkDeleteDto struct {
	Id []primitive.ObjectID `json:"id" binding:"required,dive,gt=0"`
}

func (x *Controller) BulkDelete(c *gin.Context) interface{} {
	ctx := c.Request.Context()
	var body BulkDeleteDto
	if err := c.ShouldBindJSON(&body); err != nil {
		return err
	}
	result, err := x.Service.BulkDelete(ctx, body.Id)
	if err != nil {
		return err
	}
	return result
}
