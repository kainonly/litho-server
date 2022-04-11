package pictures

import (
	"github.com/gin-gonic/gin"
)

type Controller struct {
	Service *Service
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
