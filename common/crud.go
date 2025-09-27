package common

import (
	"github.com/cloudwego/hertz/pkg/app"
	"golang.org/x/net/context"
)

type Controller interface {
	Create(ctx context.Context, c *app.RequestContext)
	Find(ctx context.Context, c *app.RequestContext)
	FindById(ctx context.Context, c *app.RequestContext)
	Update(ctx context.Context, c *app.RequestContext)
	Delete(ctx context.Context, c *app.RequestContext)
}

func SetPipe(ctx context.Context, i any) context.Context {
	return context.WithValue(ctx, "pipe", i)
}
