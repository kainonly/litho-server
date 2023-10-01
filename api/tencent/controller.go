package tencent

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
)

type Controller struct {
	TencentService *Service
}

func (x *Controller) CosPresigned(_ context.Context, c *app.RequestContext) {
	r, err := x.TencentService.CosPresigned()
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, r)
}

type CosImageInfoDto struct {
	Url string `query:"url" vd:"required"`
}

func (x *Controller) CosImageInfo(ctx context.Context, c *app.RequestContext) {
	var dto CosImageInfoDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	r, err := x.TencentService.CosImageInfo(ctx, dto.Url)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, r)
}
