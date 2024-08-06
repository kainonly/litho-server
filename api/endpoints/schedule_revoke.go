package endpoints

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	schedule "github.com/weplanx/schedule/client"
	"github.com/weplanx/server/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ScheduleRevokeDto struct {
	Id  primitive.ObjectID `json:"id" vd:"required"`
	Key string             `json:"key" vd:"required"`
}

func (x *Controller) ScheduleRevoke(ctx context.Context, c *app.RequestContext) {
	var dto ScheduleRevokeDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	if err := x.EndpointsX.ScheduleRevoke(ctx, dto.Id, dto.Key); err != nil {
		c.Error(err)
		return
	}

	c.Status(204)
}

func (x *Service) ScheduleRevoke(ctx context.Context, id primitive.ObjectID, key string) (err error) {
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
	return sc.Remove(key)
}
