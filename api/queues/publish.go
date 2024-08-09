package queues

import (
	"context"
	"github.com/bytedance/sonic"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/nats-io/nats.go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PublishDto struct {
	Project primitive.ObjectID `json:"project" vd:"required"`
	Subject string             `json:"subject" vd:"required"`
	Payload M                  `json:"payload" vd:"gt=0"`
}

func (x *Controller) Publish(ctx context.Context, c *app.RequestContext) {
	var dto PublishDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	r, err := x.QueuesX.Publish(ctx, dto)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(201, r)
}

func (x *Service) Publish(ctx context.Context, dto PublishDto) (r interface{}, err error) {
	var js nats.JetStreamContext
	if js, err = x.GetJetStream(ctx, dto.Project); err != nil {
		return
	}
	var payload []byte
	if payload, err = sonic.Marshal(dto.Payload); err != nil {
		return
	}
	if r, err = js.Publish(dto.Subject, payload, nats.Context(ctx)); err != nil {
		return
	}
	return
}
