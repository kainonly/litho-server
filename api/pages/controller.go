package pages

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/route"
)

type Controller struct {
	PagesService *Service
}

func (x *Controller) In(r *route.RouterGroup) {}

// Index 载入页面数据
func (x *Controller) Index(ctx context.Context, c *app.RequestContext) {}
