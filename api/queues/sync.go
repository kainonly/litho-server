package queues

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/nats-io/nats.go"
	"github.com/weplanx/server/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SyncDto struct {
	Id primitive.ObjectID `json:"id" vd:"mongodb"`
}

func (x *Controller) Sync(ctx context.Context, c *app.RequestContext) {
	var dto SyncDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	if err := x.QueuesX.Sync(ctx, dto.Id); err != nil {
		c.Error(err)
		return
	}

	c.Status(204)
}

func (x *Service) Sync(ctx context.Context, id primitive.ObjectID) (err error) {
	var queue model.Queue
	if err = x.Db.Collection("queues").
		FindOne(ctx, bson.M{"_id": id}).
		Decode(&queue); err != nil {
		return
	}

	var js nats.JetStreamContext
	if js, err = x.GetJetStream(ctx, queue.Project); err != nil {
		return
	}

	if _, err = js.StreamInfo(queue.ID.Hex()); err != nil {
		if err != nats.ErrStreamNotFound {
			return
		}
	}
	cfg := &nats.StreamConfig{
		Name:        queue.ID.Hex(),
		Description: queue.Name,
		Subjects:    queue.Subjects,
		MaxMsgs:     queue.MaxMsgs,
		MaxBytes:    queue.MaxBytes,
		MaxAge:      queue.MaxAge,
		Retention:   nats.WorkQueuePolicy,
	}
	if err == nats.ErrStreamNotFound {
		if _, err = js.AddStream(cfg, nats.Context(ctx)); err != nil {
			return
		}
	} else {
		if _, err = js.UpdateStream(cfg, nats.Context(ctx)); err != nil {
			return
		}
	}
	return
}
