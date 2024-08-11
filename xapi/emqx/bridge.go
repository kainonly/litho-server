package emqx

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	transfer "github.com/weplanx/collector/client"
	"time"
)

type BridgeDto struct {
	Client  string `json:"client" vd:"required"`
	Topic   string `json:"topic" vd:"required"`
	Payload M      `json:"payload" vd:"required,gt=0"`
}

func (x *Controller) Bridge(ctx context.Context, c *app.RequestContext) {
	var dto BridgeDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	if err := x.EmqxService.Bridge(ctx, dto); err != nil {
		c.Error(err)
		return
	}

	c.Status(204)
}

func (x *Service) Bridge(ctx context.Context, dto BridgeDto) (err error) {
	return x.Transfer.Publish(ctx, "logset_imessages", transfer.Payload{
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"metadata": map[string]interface{}{
				"client": dto.Client,
				"topic":  dto.Topic,
			},
			"payload": dto.Payload,
		},
	})
}
