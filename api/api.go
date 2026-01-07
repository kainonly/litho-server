package api

import (
	"context"
	"server/api/index"

	"github.com/cloudwego/hertz/pkg/app/server"
	"go.uber.org/fx"
)

var Options = fx.Options(
	index.Provides,
	fx.Invoke(Routes),
)

func Routes(lc fx.Lifecycle, h *server.Hertz, index *index.Controller) {
	h.GET("", index.Ping)

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
