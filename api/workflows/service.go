package workflows

import (
	"context"
	"github.com/bytedance/sonic"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/nats-io/nats.go"
	"github.com/weplanx/go/rest"
	sctyp "github.com/weplanx/schedule/common"
	"github.com/weplanx/server/api/endpoints"
	"github.com/weplanx/server/common"
	"github.com/weplanx/server/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Service struct {
	*common.Inject

	EndpointsX *endpoints.Service
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

func (x *Service) Event() (err error) {
	if _, err = x.JetStream.QueueSubscribe(`events.workflows`, `EVENT_workflows`, func(msg *nats.Msg) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		var dto rest.PublishDto
		if err = sonic.Unmarshal(msg.Data, &dto); err != nil {
			hlog.Error(err)
			return
		}
		switch dto.Action {
		case rest.ActionUpdateById:
			id, _ := primitive.ObjectIDFromHex(dto.Id)
			if err = x.Sync(ctx, id); err != nil {
				hlog.Error(err)
			}
			break
		case rest.ActionDelete:
			endpointId, _ := primitive.ObjectIDFromHex(dto.Data.(M)["schedule"].(M)["ref"].(string))
			if err = x.EndpointsX.ScheduleRevoke(ctx, endpointId, dto.Id); err != nil {
				hlog.Error(err)
			}
			break
		case rest.ActionBulkDelete:
			data := dto.Data.([]interface{})
			for _, v := range data {
				endpointId, _ := primitive.ObjectIDFromHex(v.(M)["schedule"].(M)["ref"].(string))
				key := v.(M)["_id"].(string)
				if err = x.EndpointsX.ScheduleRevoke(ctx, endpointId, key); err != nil {
					hlog.Error(err)
				}
			}
			break
		}
	}); err != nil {
		return
	}
	return
}
