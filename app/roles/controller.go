package roles

import "github.com/gin-gonic/gin"

type Controller struct {
	Service *Service
}

func (x *Controller) FindLabels(c *gin.Context) interface{} {
	ctx := c.Request.Context()
	result, err := x.Service.FindLabels(ctx)
	if err != nil {
		return err
	}
	return result
}

func (x *Controller) HasKey(c *gin.Context) interface{} {
	var query struct {
		Key string `form:"key" binding:"required,key"`
	}
	if err := c.ShouldBindQuery(&query); err != nil {
		return err
	}
	ctx := c.Request.Context()
	code, err := x.Service.HasKey(ctx, query.Key)
	if err != nil {
		return err
	}
	return gin.H{
		"status": code,
	}
}
