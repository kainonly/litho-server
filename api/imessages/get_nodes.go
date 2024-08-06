package imessages

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
)

func (x *Controller) GetNodes(ctx context.Context, c *app.RequestContext) {
	r, err := x.ImessagesX.GetNodes(ctx)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, r)
}

func (x *Service) GetNodes(ctx context.Context) (r interface{}, err error) {
	if _, err = x.R(ctx).
		SetSuccessResult(&r).
		SetErrorResult(&r).
		Get("nodes"); err != nil {
		return
	}
	return
}
