package app

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/route"
	"net/http"
)

type Controller struct {
	AppService *Service
}

func (x *Controller) In(r *route.RouterGroup) {
	r.GET("", x.Index)
	r.POST("auth", x.AuthLogin)
}

func (x *Controller) Index(ctx context.Context, c *app.RequestContext) {
	c.JSON(http.StatusOK, utils.H{
		"time": x.AppService.Index(),
		"ip":   c.ClientIP(),
	})
}

func (x *Controller) AuthLogin(ctx context.Context, c *app.RequestContext) {
	var dto struct {
		Identity string `json:"identity,required" vd:"len($)>=4 || email($)"`
		Password string `json:"password,required" vd:"len($)>=8"`
	}
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}
