package builders

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Controller struct {
	BuildersX *Service
}

type SortFieldsDto struct {
	Id   primitive.ObjectID `json:"id" vd:"required"`
	Keys []string           `json:"keys" vd:"gt=0"`
}

func (x *Controller) SortFields(ctx context.Context, c *app.RequestContext) {
	var dto SortFieldsDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	if err := x.BuildersX.SortFields(ctx, dto.Id, dto.Keys); err != nil {
		c.Error(err)
		return
	}

	c.Status(204)
}
