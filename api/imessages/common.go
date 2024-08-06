package imessages

import (
	"context"
	"github.com/bytedance/sonic"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/google/wire"
	"github.com/imroc/req/v3"
	"github.com/nats-io/nats.go"
	"github.com/weplanx/go/rest"
	"github.com/weplanx/server/common"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

var Provides = wire.NewSet(
	wire.Struct(new(Controller), "*"),
	wire.Struct(new(Service), "*"),
)

type Controller struct {
	ImessagesX *Service
}

type Service struct {
	*common.Inject
}

type M = map[string]interface{}

func (x *Service) R(ctx context.Context) *req.Request {
	return common.HttpClient(x.V.EmqxHost).
		SetCommonBasicAuth(x.V.EmqxApiKey, x.V.EmqxSecretKey).
		R().SetContext(ctx)
}

func (x *Service) Event() (err error) {
	if _, err = x.JetStream.QueueSubscribe(`events.imessages`, `EVENT_imessages`, func(msg *nats.Msg) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		var dto rest.PublishDto
		if err = sonic.Unmarshal(msg.Data, &dto); err != nil {
			hlog.Error(err)
			return
		}
		switch dto.Action {
		case rest.ActionCreate:
			id, _ := primitive.ObjectIDFromHex(dto.Result.(M)["InsertedID"].(string))
			if _, err = x.UpdateRule(ctx, id); err != nil {
				hlog.Error(err)
			}
			if _, err = x.UpdateMetrics(ctx, id); err != nil {
				hlog.Error(err)
			}
			break
		case rest.ActionUpdateById:
			id, _ := primitive.ObjectIDFromHex(dto.Id)
			if _, err = x.UpdateRule(ctx, id); err != nil {
				hlog.Error(err)
			}
			if _, err = x.DeleteMetrics(ctx, id); err != nil {
				hlog.Error(err)
			}
			if _, err = x.UpdateMetrics(ctx, id); err != nil {
				hlog.Error(err)
			}
			break
		}
	}); err != nil {
		return
	}
	return
}
