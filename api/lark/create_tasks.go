package lark

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
)

func (x *Controller) CreateTasks(ctx context.Context, c *app.RequestContext) {
	r, err := x.LarkX.CreateTask(ctx)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(201, r)
}

func (x *Service) CreateTask(ctx context.Context) (result M, err error) {
	var token string
	if token, err = x.GetTenantAccessToken(ctx); err != nil {
		return
	}
	body := `{}`
	if _, err = client.R().
		SetContext(ctx).
		SetBearerAuthToken(token).
		SetBody(body).
		SetSuccessResult(&result).
		Post("/task/v1/tasks"); err != nil {
		return
	}
	return
}
