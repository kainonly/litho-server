package projects

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Controller struct {
	ProjectsX *Service
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

	if err := x.ProjectsX.DeployNats(ctx, dto.Id); err != nil {
		c.Error(err)
		return
	}

	c.Status(204)
}

type GetTenantsDto struct {
	Id string `path:"id" vd:"mongodb"`
}

func (x *Controller) GetTenants(ctx context.Context, c *app.RequestContext) {
	var dto GetTenantsDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	id, _ := primitive.ObjectIDFromHex(dto.Id)
	r, err := x.ProjectsX.GetTenants(ctx, id)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, r)

}
