package schedules

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

var clients = sync.Map{}

func (x *Service) Client(node string) (client *schedule.Client, err error) {
	if i, ok := clients.Load(node); ok {
		return i.(*schedule.Client), nil
	}
	if client, err = schedule.New(node, x.Nats); err != nil {
		return
	}
	clients.Store(node, client)
	return
}

func (x *Service) Ping(node string) (r bool, err error) {
	var sc *schedule.Client
	if sc, err = x.Client(node); err != nil {
		return
	}
	return sc.Ping()
}

func (x *Service) Keys(ctx context.Context, id primitive.ObjectID) (keys []string, err error) {
	var data model.Schedule
	if err = x.Db.Collection("schedules").
		FindOne(ctx, bson.M{"_id": id}).
		Decode(&data); err != nil {
		return
	}
	var sc *schedule.Client
	if sc, err = x.Client(data.Node); err != nil {
		return
	}
	return sc.Lists()
}

func (x *Service) Set(ctx context.Context, id primitive.ObjectID, key string, option sctyp.ScheduleOption) (err error) {
	var data model.Schedule
	if err = x.Db.Collection("schedules").
		FindOne(ctx, bson.M{"_id": id}).
		Decode(&data); err != nil {
		return
	}
	var sc *schedule.Client
	if sc, err = x.Client(data.Node); err != nil {
		return
	}
	return sc.Set(key, option)
}

func (x *Service) Revoke(ctx context.Context, id primitive.ObjectID, key string) (err error) {
	var data model.Schedule
	if err = x.Db.Collection("schedules").
		FindOne(ctx, bson.M{"_id": id}).
		Decode(&data); err != nil {
		return
	}
	var sc *schedule.Client
	if sc, err = x.Client(data.Node); err != nil {
		return
	}
	return sc.Remove(key)
}

func (x *Service) State(node string, key string) (r sctyp.ScheduleOption, err error) {
	var sc *schedule.Client
	if sc, err = x.Client(node); err != nil {
		return
	}
	return sc.Get(key)
}
