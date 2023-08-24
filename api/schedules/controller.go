package schedules

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

type Controller struct {
	SchedulesService *Service
}

type PingDto struct {
	Id string `path:"id,required"`
}

func (x *Controller) Ping(ctx context.Context, c *app.RequestContext) {
	var dto PingDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	if _, err := x.SchedulesService.Ping(dto.Id); err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusOK)
}

type DeployDto struct {
	Id string `path:"id,required"`
}

func (x *Controller) Deploy(ctx context.Context, c *app.RequestContext) {
	var dto DeployDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	id, _ := primitive.ObjectIDFromHex(dto.Id)
	if err := x.SchedulesService.Deploy(ctx, id); err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusOK)
}

type UnDeployDto struct {
	Id string `path:"id,required"`
}

func (x *Controller) Undeploy(ctx context.Context, c *app.RequestContext) {
	var dto UnDeployDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	id, _ := primitive.ObjectIDFromHex(dto.Id)
	if err := x.SchedulesService.Undeploy(ctx, id); err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusOK)
}
