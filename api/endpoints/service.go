package endpoints

import (
	"context"
	schedule "github.com/weplanx/schedule/client"
	sctyp "github.com/weplanx/schedule/common"
	"github.com/weplanx/server/common"
	"github.com/weplanx/server/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"sync"
)

type Service struct {
	*common.Inject
}

var schedules = sync.Map{}

func (x *Service) Schedule(node string) (client *schedule.Client, err error) {
	if i, ok := schedules.Load(node); ok {
		return i.(*schedule.Client), nil
	}
	if client, err = schedule.New(node, x.Nats); err != nil {
		return
	}
	schedules.Store(node, client)
	return
}

func (x *Service) SchedulePing(node string) (r bool, err error) {
	var sc *schedule.Client
	if sc, err = x.Schedule(node); err != nil {
		return
	}
	return sc.Ping()
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

func (x *Service) ScheduleSet(ctx context.Context, id primitive.ObjectID, key string, option sctyp.ScheduleOption) (err error) {
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
	return sc.Set(key, option)
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

func (x *Service) ScheduleState(node string, key string) (r sctyp.ScheduleOption, err error) {
	var sc *schedule.Client
	if sc, err = x.Schedule(node); err != nil {
		return
	}
	return sc.Get(key)
}
