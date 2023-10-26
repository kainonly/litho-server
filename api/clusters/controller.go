package clusters

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Controller struct {
	ClustersX *Service
}

type GetInfoDto struct {
	Id string `path:"id" vd:"mongodb"`
}

func (x *Controller) GetInfo(ctx context.Context, c *app.RequestContext) {
	var dto GetInfoDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	id, _ := primitive.ObjectIDFromHex(dto.Id)
	r, err := x.ClustersX.GetInfo(ctx, id)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, r)
}

type GetNodesDto struct {
	Id string `path:"id" vd:"mongodb"`
}

func (x *Controller) GetNodes(ctx context.Context, c *app.RequestContext) {
	var dto GetNodesDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	id, _ := primitive.ObjectIDFromHex(dto.Id)
	r, err := x.ClustersX.GetNodes(ctx, id)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, r)
}
