package schedules

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Controller struct {
	Service *Service
}

// List 获取配置唯一标识
func (x *Controller) List(c *gin.Context) interface{} {
	keys, err := x.Service.List()
	if err != nil {
		return err
	}
	return keys
}

// Get 获取指定服务配置与运行状态
func (x *Controller) Get(c *gin.Context) interface{} {
	var uri struct {
		Key string `uri:"key"`
	}
	if err := c.ShouldBindUri(&uri); err != nil {
		return err
	}
	result, err := x.Service.Get(uri.Key)
	if err != nil {
		return err
	}
	return result
}

// SetSync 设置同步
func (x *Controller) SetSync(c *gin.Context) interface{} {
	var body struct {
		Id primitive.ObjectID `json:"id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		return err
	}
	ctx := c.Request.Context()
	if err := x.Service.Sync(ctx, body.Id); err != nil {
		return err
	}
	return nil
}
