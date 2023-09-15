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
	Nodes []string `json:"nodes,required"`
}

func (x *Controller) Ping(_ context.Context, c *app.RequestContext) {
	var dto PingDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	result := make(M)
	for _, node := range dto.Nodes {
		r, err := x.SchedulesService.Ping(node)
		if err != nil {
			c.Error(err)
			return
		}
		result[node] = r
	}

	c.JSON(http.StatusOK, result)
}

type KeysDto struct {
	Id string `path:"id,required" vd:"mongoId($);msg:'the document id must be an ObjectId'"`
}

func (x *Controller) Keys(ctx context.Context, c *app.RequestContext) {
	var dto KeysDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	id, _ := primitive.ObjectIDFromHex(dto.Id)
	r, err := x.SchedulesService.Keys(ctx, id)
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

type StatesDto struct {
	Node string `json:"node,required"`
	Key  string `json:"key,required"`
}

func (x *Controller) State(_ context.Context, c *app.RequestContext) {
	var dto StatesDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	r, err := x.SchedulesService.State(dto.Node, dto.Key)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, r)
}
