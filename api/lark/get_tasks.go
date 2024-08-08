package lark

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
)

func (x *Controller) GetTasks(ctx context.Context, c *app.RequestContext) {
	r, err := x.LarkX.GetTasks(ctx)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, r)
}

func (x *Service) GetTasks(ctx context.Context) (result M, err error) {
	var token string
	if token, err = x.GetTenantAccessToken(ctx); err != nil {
		return
	}
	if _, err = client.R().
		SetContext(ctx).
		SetBearerAuthToken(token).
		SetSuccessResult(&result).
		Get("/task/v1/tasks"); err != nil {
		return
	}
	return
}
