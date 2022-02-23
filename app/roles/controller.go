package roles

import "github.com/gin-gonic/gin"

type Controller struct {
	Service *Service
}

func (x *Controller) HasName(c *gin.Context) interface{} {
	var query struct {
		Name string `form:"name" binding:"required"`
	}
	if err := c.ShouldBindQuery(&query); err != nil {
		return err
	}
	ctx := c.Request.Context()
	code, err := x.Service.HasName(ctx, query.Name)
	if err != nil {
		return err
	}
	return gin.H{
		"status": code,
	}
}

func (x *Controller) Labels(c *gin.Context) interface{} {
	ctx := c.Request.Context()
	values, err := x.Service.FindLabels(ctx)
	if err != nil {
		return err
	}
	return values
}
