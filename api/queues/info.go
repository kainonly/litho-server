package queues

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/nats-io/nats.go"
	"github.com/weplanx/server/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StateDto struct {
	Id string `path:"id" vd:"mongodb"`
}

func (x *Controller) Info(ctx context.Context, c *app.RequestContext) {
	var dto StateDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	id, _ := primitive.ObjectIDFromHex(dto.Id)
	r, err := x.QueuesX.Info(ctx, id)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, r)
}

func (x *Service) Info(ctx context.Context, id primitive.ObjectID) (r *nats.StreamInfo, err error) {
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

	if r, err = js.StreamInfo(queue.ID.Hex(), nats.Context(ctx)); err != nil {
		return
	}
	r.Cluster = nil
	return
}
