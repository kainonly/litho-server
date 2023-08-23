package schedules

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"net/http"
)

type Controller struct {
	SchedulesService *Service
}

type PingDto struct {
	NodeId string `path:"node_id,required"`
}

func (x *Controller) Ping(ctx context.Context, c *app.RequestContext) {
	var dto PingDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	if _, err := x.SchedulesService.Ping(dto.NodeId); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, utils.H{
		"msg": "ok",
	})
}
