package endpoints

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	schedule "github.com/weplanx/schedule/client"
)

type SchedulePingDto struct {
	Nodes []string `json:"nodes" vd:"gt=0"`
}

func (x *Controller) SchedulePing(_ context.Context, c *app.RequestContext) {
	var dto SchedulePingDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	result := make(M)
	for _, node := range dto.Nodes {
		r, err := x.EndpointsX.SchedulePing(node)
		if err != nil {
			c.Error(err)
			return
		}
		result[node] = r
	}

	c.JSON(200, result)
}

func (x *Service) SchedulePing(node string) (r bool, err error) {
	var sc *schedule.Client
	if sc, err = x.Schedule(node); err != nil {
		return
	}
	return sc.Ping()
}
