package main

import (
	"context"
	"server/api"
	"server/bootstrap"

	"github.com/cloudwego/hertz/pkg/app/server"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		bootstrap.Provides,
		api.Options,
		fx.Invoke(NewServer),
	).Run()
}

func NewServer(lc fx.Lifecycle, api *api.API) (err error) {
	ctx := context.TODO()
	var h *server.Hertz
	if h, err = api.Initialize(ctx); err != nil {
		return
	}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go h.Spin()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			if h != nil {
				return h.Shutdown(ctx)
			}
			return nil
		},
	})
	return
}
