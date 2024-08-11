package index

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/weplanx/server/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (x *Controller) Accelerate(ctx context.Context, c *app.RequestContext) {
	r, err := x.IndexService.Accelerate(ctx)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, r)
}

func (x *Service) Accelerate(ctx context.Context) (result []model.AccTask, err error) {
	var cursor *mongo.Cursor
	if cursor, err = x.Db.Collection("acc_tasks").Find(ctx, bson.M{}); err != nil {
		return
	}
	if err = cursor.All(ctx, &result); err != nil {
		return
	}
	return
}
