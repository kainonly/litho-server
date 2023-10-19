package acc_tasks

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
)

type Controller struct {
	AccTasksService *Service
}

func (x *Controller) Invoke(ctx context.Context, c *app.RequestContext) {
	r, err := x.AccTasksService.Invoke(ctx)

	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, r)
}
