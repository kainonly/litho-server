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
	Ids []string `json:"ids,required"`
}

func (x *Controller) Ping(_ context.Context, c *app.RequestContext) {
	var dto PingDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	result := make(M)
	for _, id := range dto.Ids {
		r, err := x.SchedulesService.Ping(id)
		if err != nil {
			c.Error(err)
			return
		}
		result[id] = r
	}

	c.JSON(http.StatusOK, result)
}

type DeployDto struct {
	Id primitive.ObjectID `json:"id,required"`
}

func (x *Controller) Deploy(ctx context.Context, c *app.RequestContext) {
	var dto DeployDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	if err := x.SchedulesService.Deploy(ctx, dto.Id); err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusOK)
}

type UnDeployDto struct {
	Id primitive.ObjectID `json:"id,required"`
}

func (x *Controller) Undeploy(ctx context.Context, c *app.RequestContext) {
	var dto UnDeployDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	if err := x.SchedulesService.Undeploy(ctx, dto.Id); err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusOK)
}

type KeysDto struct {
	Id string `path:"id,required"`
}

func (x *Controller) Keys(_ context.Context, c *app.RequestContext) {
	var dto KeysDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	r, err := x.SchedulesService.Keys(dto.Id)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, r)
}

type RevokeDto struct {
	Id  string `json:"id,required"`
	Key string `json:"key,required"`
}

func (x *Controller) Revoke(_ context.Context, c *app.RequestContext) {
	var dto RevokeDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	if err := x.SchedulesService.Revoke(dto.Id, dto.Key); err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusOK)
}
