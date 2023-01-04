package tencent

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"net/http"
)

type Controller struct {
	TencentService *Service
}

// CosPresigned 对象存储预签名
func (x *Controller) CosPresigned(ctx context.Context, c *app.RequestContext) {
	result, err := x.TencentService.CosPresigned(ctx)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, result)
}

type ImageInfoDto struct {
	Url string `query:"url,required"`
}

// ImageInfo 获取图片信息
func (x *Controller) ImageInfo(ctx context.Context, c *app.RequestContext) {
	var dto ImageInfoDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}
	result, err := x.TencentService.ImageInfo(ctx, dto.Url)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, result)
}
