package endpoints

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Controller struct {
	EndpointsX *Service
}

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

type ScheduleKeysDto struct {
	Id string `path:"id" vd:"mongodb"`
}

func (x *Controller) ScheduleKeys(ctx context.Context, c *app.RequestContext) {
	var dto ScheduleKeysDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	id, _ := primitive.ObjectIDFromHex(dto.Id)
	r, err := x.EndpointsX.ScheduleKeys(ctx, id)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, r)
}

type ScheduleRevokeDto struct {
	Id  primitive.ObjectID `json:"id" vd:"required"`
	Key string             `json:"key" vd:"required"`
}

func (x *Controller) ScheduleRevoke(ctx context.Context, c *app.RequestContext) {
	var dto ScheduleRevokeDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	if err := x.EndpointsX.ScheduleRevoke(ctx, dto.Id, dto.Key); err != nil {
		c.Error(err)
		return
	}

	c.Status(204)
}

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
