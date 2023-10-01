package imessages

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Controller struct {
	ImessagesServices *Service
}

func (x *Controller) GetNodes(ctx context.Context, c *app.RequestContext) {
	r, err := x.ImessagesServices.GetNodes(ctx)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, r)
}

type GetMetricsDto struct {
	Id string `path:"id" vd:"mongodb"`
}

func (x *Controller) GetMetrics(ctx context.Context, c *app.RequestContext) {
	var dto GetMetricsDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	id, _ := primitive.ObjectIDFromHex(dto.Id)
	r, err := x.ImessagesServices.GetMetrics(ctx, id)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, r)
}

type CreateMetricsDto struct {
	Id string `path:"id" vd:"mongodb"`
}

func (x *Controller) CreateMetrics(ctx context.Context, c *app.RequestContext) {
	var dto CreateMetricsDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	id, _ := primitive.ObjectIDFromHex(dto.Id)
	r, err := x.ImessagesServices.CreateMetrics(ctx, id)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(201, r)
}

type DeleteMetricsDto struct {
	Id string `path:"id" vd:"mongodb"`
}

func (x *Controller) DeleteMetrics(ctx context.Context, c *app.RequestContext) {
	var dto DeleteMetricsDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	id, _ := primitive.ObjectIDFromHex(dto.Id)
	r, err := x.ImessagesServices.DeleteMetrics(ctx, id)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, r)
}

type PublishDto struct {
	Topic   string `json:"topic" vd:"required"`
	Payload M      `json:"payload" vd:"required"`
}

func (x *Controller) Publish(ctx context.Context, c *app.RequestContext) {
	var dto PublishDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	r, err := x.ImessagesServices.Publish(ctx, dto)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(201, r)
}
