package datasets

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
)

type DeleteDto struct {
	Name string `path:"name" vd:"required"`
}

func (x *Controller) Delete(ctx context.Context, c *app.RequestContext) {
	var dto DeleteDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	if err := x.DatasetsX.Delete(ctx, dto.Name); err != nil {
		c.Error(err)
		return
	}

	c.Status(204)
}

func (x *Service) Delete(ctx context.Context, name string) (err error) {
	if err = x.Db.Collection(name).Drop(ctx); err != nil {
		return
	}
	controls := x.V.RestControls
	delete(controls, name)
	if err = x.Values.Set(M{
		"RestControls": controls,
	}); err != nil {
		return
	}
	return
}
