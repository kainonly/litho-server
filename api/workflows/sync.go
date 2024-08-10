package workflows

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	sctyp "github.com/weplanx/schedule/common"
	"github.com/weplanx/server/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SyncDto struct {
	Id primitive.ObjectID `json:"id" vd:"required"`
}

func (x *Controller) Sync(ctx context.Context, c *app.RequestContext) {
	var dto SyncDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	if err := x.WorkflowsX.Sync(ctx, dto.Id); err != nil {
		c.Error(err)
		return
	}

	c.Status(200)
}

func (x *Service) Sync(ctx context.Context, id primitive.ObjectID) (err error) {
	var data model.Workflow
	if err = x.Db.Collection("workflows").FindOne(ctx, bson.M{
		"_id": id,
	}).Decode(&data); err != nil {
		return
	}

	if data.Schedule != nil {
		schedule := data.Schedule
		jobs := make([]sctyp.ScheduleJob, len(schedule.Jobs))
		for i, v := range schedule.Jobs {
			jobs[i] = sctyp.ScheduleJob{
				Mode:   v.Mode,
				Spec:   v.Spec,
				Option: v.Option,
			}
		}
		if err = x.EndpointsX.ScheduleSet(ctx,
			*schedule.Ref,
			id.Hex(),
			sctyp.ScheduleOption{
				Status: schedule.Status,
				Jobs:   jobs,
			},
		); err != nil {
			return
		}
	}
	return
}
