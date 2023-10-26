package imessages

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Controller struct {
	ImessagesX *Service
}

func (x *Controller) GetNodes(ctx context.Context, c *app.RequestContext) {
	r, err := x.ImessagesX.GetNodes(ctx)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, r)
}

type UpdateRuleDto struct {
	Id string `path:"id" vd:"mongodb"`
}

func (x *Controller) UpdateRule(ctx context.Context, c *app.RequestContext) {
	var dto UpdateRuleDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	id, _ := primitive.ObjectIDFromHex(dto.Id)
	r, err := x.ImessagesX.UpdateRule(ctx, id)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(201, r)
}

type DeleteRuleDto struct {
	Id string `path:"id" vd:"mongodb"`
}

func (x *Controller) DeleteRule(ctx context.Context, c *app.RequestContext) {
	var dto DeleteRuleDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	id, _ := primitive.ObjectIDFromHex(dto.Id)
	if err := x.ImessagesX.DeleteRule(ctx, id); err != nil {
		c.Error(err)
		return
	}

	c.Status(204)
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
	r, err := x.ImessagesX.GetMetrics(ctx, id)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, r)
}

type UpdateMetricsDto struct {
	Id string `path:"id" vd:"mongodb"`
}

func (x *Controller) UpdateMetrics(ctx context.Context, c *app.RequestContext) {
	var dto UpdateMetricsDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	id, _ := primitive.ObjectIDFromHex(dto.Id)
	r, err := x.ImessagesX.UpdateMetrics(ctx, id)
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
	r, err := x.ImessagesX.DeleteMetrics(ctx, id)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, r)
}

type PublishDto struct {
	Topic   string `json:"topic" vd:"required"`
	Payload M      `json:"payload" vd:"required,gt=0"`
}

func (x *Controller) Publish(ctx context.Context, c *app.RequestContext) {
	var dto PublishDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	r, err := x.ImessagesX.Publish(ctx, dto)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(201, r)
}
