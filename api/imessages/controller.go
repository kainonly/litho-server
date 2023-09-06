package imessages

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
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

	c.JSON(http.StatusOK, r)
}

type GetMetricsDto struct {
	Id string `path:"id,required" vd:"mongoId($);msg:'the document id must be an ObjectId'"`
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

	c.JSON(http.StatusOK, r)
}

type CreateMetricsDto struct {
	Id string `path:"id,required" vd:"mongoId($);msg:'the document id must be an ObjectId'"`
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

	c.JSON(http.StatusOK, r)
}

type DeleteMetricsDto struct {
	Id string `path:"id,required" vd:"mongoId($);msg:'the document id must be an ObjectId'"`
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

	c.JSON(http.StatusOK, r)
}

type PublishDto struct {
	Topic   string `json:"topic,required""`
	Payload M      `json:"payload,required"`
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

	c.JSON(http.StatusOK, r)
}
