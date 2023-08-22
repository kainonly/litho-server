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
