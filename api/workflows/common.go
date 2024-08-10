package workflows

import (
	"context"
	"github.com/bytedance/sonic"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/google/wire"
	"github.com/nats-io/nats.go"
	"github.com/weplanx/go/rest"
	"github.com/weplanx/server/api/endpoints"
	"github.com/weplanx/server/common"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

var Provides = wire.NewSet(
	wire.Struct(new(Controller), "*"),
	wire.Struct(new(Service), "*"),
)

type Controller struct {
	WorkflowsX *Service
}

type Service struct {
	*common.Inject

	EndpointsX *endpoints.Service
}

type M = map[string]interface{}

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
