package values

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/route"
	"net/http"
)

type Controller struct {
	ValuesService *Service
}

func (x *Controller) In(r *route.RouterGroup) {
	r.GET("", x.Get)
	r.PATCH("", x.Set)
	r.DELETE(":key", x.Remove)
}

// Get 获取动态配置
func (x *Controller) Get(ctx context.Context, c *app.RequestContext) {
	var query struct {
		// 动态配置键
		Keys []string `query:"keys"`
	}
	if err := c.BindAndValidate(&query); err != nil {
		c.Error(err)
		return
	}

	data := x.ValuesService.Get(query.Keys...)

	c.JSON(http.StatusOK, utils.H{"data": data})
}

// Set 设置动态配置
func (x *Controller) Set(ctx context.Context, c *app.RequestContext) {
	var body struct {
		Data map[string]interface{} `json:"data,required" vd:"len($)>0"`
	}
	if err := c.BindAndValidate(&body); err != nil {
		c.Error(err)
		return
	}

	if err := x.ValuesService.Set(ctx, body.Data); err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}

// Remove 移除动态配置
func (x *Controller) Remove(ctx context.Context, c *app.RequestContext) {
	var params struct {
		Key string `path:"key,required"`
	}
	if err := c.BindAndValidate(&params); err != nil {
		c.Error(err)
		return
	}

	if err := x.ValuesService.Remove(ctx, params.Key); err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}
