package workflows

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

type Controller struct {
	WorkflowsService *Service
}

type SyncDto struct {
	Id primitive.ObjectID `json:"id,required"`
}

func (x *Controller) Sync(ctx context.Context, c *app.RequestContext) {
	var dto SyncDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	if err := x.WorkflowsService.Sync(ctx, dto.Id); err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusOK)
}
