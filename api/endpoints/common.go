package endpoints

import (
	"context"
	"github.com/google/wire"
	schedule "github.com/weplanx/schedule/client"
	sctyp "github.com/weplanx/schedule/common"
	"github.com/weplanx/server/common"
	"github.com/weplanx/server/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"sync"
)

var Provides = wire.NewSet(
	wire.Struct(new(Controller), "*"),
	wire.Struct(new(Service), "*"),
)

type Controller struct {
	EndpointsX *Service
}

type Service struct {
	*common.Inject
}

type M = map[string]interface{}

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
