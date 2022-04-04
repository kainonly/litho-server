package departments

import (
	"github.com/gin-gonic/gin"
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
