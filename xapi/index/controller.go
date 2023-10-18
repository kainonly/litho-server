package index

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
)

type Controller struct {
	IndexService *Service
}

func (x *Controller) Accelerate(ctx context.Context, c *app.RequestContext) {
	r, err := x.IndexService.Accelerate(ctx)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, r)
}
