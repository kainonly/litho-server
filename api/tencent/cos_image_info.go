package tencent

import (
	"context"
	"github.com/bytedance/sonic/decoder"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/tencentyun/cos-go-sdk-v5"
)

type CosImageInfoDto struct {
	Url string `query:"url" vd:"required"`
}

func (x *Controller) CosImageInfo(ctx context.Context, c *app.RequestContext) {
	var dto CosImageInfoDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	r, err := x.TencentX.CosImageInfo(ctx, dto.Url)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, r)
}

func (x *Service) CosImageInfo(ctx context.Context, url string) (r M, err error) {
	client := x.Cos()
	var res *cos.Response
	if res, err = client.CI.Get(ctx, url, "imageInfo", nil); err != nil {
		return
	}
	if err = decoder.NewStreamDecoder(res.Body).Decode(&r); err != nil {
		return
	}
	return
}
