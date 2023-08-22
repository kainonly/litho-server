package datasets

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"net/http"
)

type Controller struct {
	DatasetsService *Service
}

type ListsDto struct {
	Name string `query:"name,required"`
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

	c.JSON(http.StatusOK, r)
}

type CreateDto struct {
	Name   string          `json:"name"`
	Kind   string          `json:"kind"`
	Option CreateOptionDto `json:"option,omitempty"`
}

type CreateOptionDto struct {
	Time   string `json:"time"`
	Meta   string `json:"meta"`
	Expire int64  `json:"expire"`
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

	c.Status(http.StatusOK)
}

type DeleteDto struct {
	Name string `path:"name,required"`
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

	c.Status(http.StatusOK)
}
