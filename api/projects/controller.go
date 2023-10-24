package projects

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Controller struct {
	ProjectsServices *Service
}

type DeployNatsDto struct {
	Id primitive.ObjectID `json:"id" vd:"required"`
}

func (x *Controller) DeployNats(ctx context.Context, c *app.RequestContext) {
	var dto DeployNatsDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	if err := x.ProjectsServices.DeployNats(ctx, dto.Id); err != nil {
		c.Error(err)
		return
	}

	c.Status(204)
}
