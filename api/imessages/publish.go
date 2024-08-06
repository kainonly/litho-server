package imessages

import (
	"context"
	"github.com/bytedance/sonic"
	"github.com/cloudwego/hertz/pkg/app"
)

type PublishDto struct {
	Topic   string `json:"topic" vd:"required"`
	Payload M      `json:"payload" vd:"required,gt=0"`
}

func (x *Controller) Publish(ctx context.Context, c *app.RequestContext) {
	var dto PublishDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	r, err := x.ImessagesX.Publish(ctx, dto)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(201, r)
}

func (x *Service) Publish(ctx context.Context, dto PublishDto) (r interface{}, err error) {
	var payload string
	if payload, err = sonic.MarshalString(dto.Payload); err != nil {
		return
	}
	if _, err = x.R(ctx).
		SetBody(M{
			"topic":   dto.Topic,
			"payload": payload,
		}).
		SetSuccessResult(&r).
		SetErrorResult(&r).
		Post("publish"); err != nil {
		return
	}
	return
}
