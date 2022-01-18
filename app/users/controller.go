package users

import "github.com/gin-gonic/gin"

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

func (x *Controller) HasUsername(c *gin.Context) interface{} {
	var query struct {
		Username string `form:"username" binding:"required,username"`
	}
	if err := c.ShouldBindQuery(&query); err != nil {
		return err
	}
	ctx := c.Request.Context()
	code, err := x.Service.HasUsername(ctx, query.Username)
	if err != nil {
		return err
	}
	return gin.H{
		"status": code,
	}
}
