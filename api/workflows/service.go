package workflows

import (
	"context"
	"github.com/bytedance/sonic"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/nats-io/nats.go"
	"github.com/weplanx/go/rest"
	"github.com/weplanx/server/api/schedules"
	"github.com/weplanx/server/common"
	"github.com/weplanx/server/model"
	"github.com/weplanx/workflow/typ"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type Service struct {
	*common.Inject

	Schedules *schedules.Service
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
		jobs := make([]typ.ScheduleJob, len(schedule.Jobs))
		for i, v := range schedule.Jobs {
			jobs[i] = typ.ScheduleJob{
				Mode:   v.Mode,
				Spec:   v.Spec,
				Option: v.Option,
			}
		}
		if err = x.Schedules.Set(
			schedule.ScheduleId.Hex(),
			id.Hex(),
			typ.ScheduleOption{
				Status: schedule.Status,
				Jobs:   jobs,
			},
		); err != nil {
			return
		}
	}
	return
}

func (x *Service) States(ctx context.Context, ids []primitive.ObjectID) (r M, err error) {
	var cursor *mongo.Cursor
	if cursor, err = x.Db.Collection("workflows").Find(ctx, bson.M{
		"_id":      bson.M{"$in": ids},
		"schedule": bson.M{"$exists": true},
	}); err != nil {
		return
	}
	r = make(M)
	for cursor.Next(ctx) {
		var data model.Workflow
		if err = cursor.Decode(&data); err != nil {
			return
		}
		var opt typ.ScheduleOption
		if opt, err = x.Schedules.Get(
			data.Schedule.ScheduleId.Hex(),
			data.ID.Hex(),
		); err != nil {
			if err == nats.ErrKeyNotFound {
				err = nil
				continue
			}
			return
		}
		r[data.ID.Hex()] = opt
	}
	return
}

func (x *Service) Event() (err error) {
	subj := x.V.NameX(".", "events", "workflows")
	queue := x.V.Name("events", "workflows")
	if _, err = x.JetStream.QueueSubscribe(subj, queue, func(msg *nats.Msg) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		var dto rest.PublishDto
		if err = sonic.Unmarshal(msg.Data, &dto); err != nil {
			hlog.Error(err)
			return
		}
		switch dto.Action {
		case "update_by_id":
			id, _ := primitive.ObjectIDFromHex(dto.Id)
			if err = x.Sync(ctx, id); err != nil {
				hlog.Error(err)
			}
			break
		}
	}); err != nil {
		return
	}
	return
}
