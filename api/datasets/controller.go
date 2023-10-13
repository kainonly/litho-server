package datasets

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
)

type Controller struct {
	DatasetsService *Service
}

type ListsDto struct {
	Name string `query:"name"`
}

func (x *Controller) Lists(ctx context.Context, c *app.RequestContext) {
	var dto ListsDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	r, err := x.DatasetsService.Lists(ctx, dto.Name)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, r)
}

type CreateDto struct {
	Name   string          `json:"name" vd:"required"`
	Kind   string          `json:"kind" vd:"oneof='collection' 'timeseries'"`
	Option CreateOptionDto `json:"option" vd:"required_if=Kind timeseries"`
}

type CreateOptionDto struct {
	Time   string `json:"time" vd:"required"`
	Meta   string `json:"meta" vd:"required"`
	Expire *int64 `json:"expire" vd:"required"`
}

func (x *Controller) Create(ctx context.Context, c *app.RequestContext) {
	var dto CreateDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	if err := x.DatasetsService.Create(ctx, dto); err != nil {
		c.Error(err)
		return
	}

	c.Status(201)
}

type DeleteDto struct {
	Name string `path:"name" vd:"required"`
}

func (x *Controller) Delete(ctx context.Context, c *app.RequestContext) {
	var dto DeleteDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	if err := x.DatasetsService.Delete(ctx, dto.Name); err != nil {
		c.Error(err)
		return
	}

	c.Status(204)
}
