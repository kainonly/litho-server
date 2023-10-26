package schedules

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Controller struct {
	SchedulesX *Service
}

type PingDto struct {
	Nodes []string `json:"nodes" vd:"gt=0"`
}

func (x *Controller) Ping(_ context.Context, c *app.RequestContext) {
	var dto PingDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	result := make(M)
	for _, node := range dto.Nodes {
		r, err := x.SchedulesX.Ping(node)
		if err != nil {
			c.Error(err)
			return
		}
		result[node] = r
	}

	c.JSON(200, result)
}

type KeysDto struct {
	Id string `path:"id" vd:"mongodb"`
}

func (x *Controller) Keys(ctx context.Context, c *app.RequestContext) {
	var dto KeysDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	id, _ := primitive.ObjectIDFromHex(dto.Id)
	r, err := x.SchedulesX.Keys(ctx, id)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, r)
}

type RevokeDto struct {
	Id  primitive.ObjectID `json:"id" vd:"required"`
	Key string             `json:"key" vd:"required"`
}

func (x *Controller) Revoke(ctx context.Context, c *app.RequestContext) {
	var dto RevokeDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	if err := x.SchedulesX.Revoke(ctx, dto.Id, dto.Key); err != nil {
		c.Error(err)
		return
	}

	c.Status(204)
}

type StatesDto struct {
	Node string `json:"node" vd:"required"`
	Key  string `json:"key" vd:"required"`
}

func (x *Controller) State(_ context.Context, c *app.RequestContext) {
	var dto StatesDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	r, err := x.SchedulesX.State(dto.Node, dto.Key)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, r)
}
