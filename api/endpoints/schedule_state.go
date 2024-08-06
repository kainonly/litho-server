package endpoints

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	schedule "github.com/weplanx/schedule/client"
	sctyp "github.com/weplanx/schedule/common"
)

type ScheduleStatesDto struct {
	Node string `json:"node" vd:"required"`
	Key  string `json:"key" vd:"required"`
}

func (x *Controller) ScheduleState(_ context.Context, c *app.RequestContext) {
	var dto ScheduleStatesDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	r, err := x.EndpointsX.ScheduleState(dto.Node, dto.Key)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, r)
}

func (x *Service) ScheduleState(node string, key string) (r sctyp.ScheduleOption, err error) {
	var sc *schedule.Client
	if sc, err = x.Schedule(node); err != nil {
		return
	}
	return sc.Get(key)
}
