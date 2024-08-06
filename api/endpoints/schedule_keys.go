package endpoints

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	schedule "github.com/weplanx/schedule/client"
	"github.com/weplanx/server/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ScheduleKeysDto struct {
	Id string `path:"id" vd:"mongodb"`
}

func (x *Controller) ScheduleKeys(ctx context.Context, c *app.RequestContext) {
	var dto ScheduleKeysDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	id, _ := primitive.ObjectIDFromHex(dto.Id)
	r, err := x.EndpointsX.ScheduleKeys(ctx, id)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, r)
}

func (x *Service) ScheduleKeys(ctx context.Context, id primitive.ObjectID) (keys []string, err error) {
	var data model.Endpoint
	if err = x.Db.Collection("endpoints").
		FindOne(ctx, bson.M{"_id": id}).
		Decode(&data); err != nil {
		return
	}
	var sc *schedule.Client
	if sc, err = x.Schedule(data.Schedule.Node); err != nil {
		return
	}
	return sc.Lists()
}
