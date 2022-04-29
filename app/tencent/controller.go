package tencent

import "github.com/gin-gonic/gin"

type Controller struct {
	Tencent *Service
}

// CosPresigned 对象存储预签名
func (x *Controller) CosPresigned(c *gin.Context) interface{} {
	data, err := x.Tencent.CosPresigned(c.Request.Context())
	if err != nil {
		return err
	}
	return data
}

func (x *Controller) ImageInfo(c *gin.Context) interface{} {
	var params struct {
		Url string `form:"url" binding:"required"`
	}
	if err := c.ShouldBindQuery(&params); err != nil {
		return err
	}
	ctx := c.Request.Context()
	result, err := x.Tencent.ImageInfo(ctx, params.Url)
	if err != nil {
		return err
	}
	return result
}
