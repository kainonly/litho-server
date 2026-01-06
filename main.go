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
	).Run()

}

func RegisterServer(lc fx.Lifecycle, h *server.Hertz) {
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
}
