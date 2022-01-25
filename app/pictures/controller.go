package pictures

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

func (x *Controller) BulkDelete(c *gin.Context) interface{} {
	var body struct {
		Id []primitive.ObjectID `json:"id" binding:"required,dive,gt=0"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		return err
	}
	ctx := c.Request.Context()
	result, err := x.Service.BulkDelete(ctx, body.Id)
	if err != nil {
		return err
	}
	return result
}

func (x *Controller) ImageInfo(c *gin.Context) interface{} {
	var params struct {
		Url string `form:"url" binding:"required"`
	}
	if err := c.ShouldBindQuery(&params); err != nil {
		return err
	}
	ctx := c.Request.Context()
	result, err := x.Service.ImageInfo(ctx, params.Url)
	if err != nil {
		return err
	}
	return result
}
