package monitor

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"net/http"
)

type Controller struct {
	MonitorService *Service
}

func (x *Controller) GetCgoCalls(ctx context.Context, c *app.RequestContext) {
	data, err := x.MonitorService.GetCgoCalls(ctx)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, data)
}
