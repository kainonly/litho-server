package clusters

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

type Controller struct {
	ClustersService *Service
}

type GetInfoDto struct {
	id string `path:"id,required" vd:"mongoId($);msg:'the document id must be an ObjectId'"`
}

func (x *Controller) GetInfo(ctx context.Context, c *app.RequestContext) {
	var dto GetInfoDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	id, _ := primitive.ObjectIDFromHex(dto.id)
	r, err := x.ClustersService.GetInfo(ctx, id)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, r)
}

type GetNodesDto struct {
	id string `path:"id,required" vd:"mongoId($);msg:'the document id must be an ObjectId'"`
}

func (x *Controller) GetNodes(ctx context.Context, c *app.RequestContext) {
	var dto GetNodesDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	id, _ := primitive.ObjectIDFromHex(dto.id)
	r, err := x.ClustersService.GetNodes(ctx, id)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, r)
}
