package xapi

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/google/wire"
	transfer "github.com/weplanx/collector/client"
	"github.com/weplanx/server/common"
	"github.com/weplanx/server/xapi/emqx"
)

var Provides = wire.NewSet(
	emqx.Provides,
)

type API struct {
	*common.Inject

	Hertz       *server.Hertz
	Emqx        *emqx.Controller
	EmqxService *emqx.Service
}

func (x *API) Routes(h *server.Hertz) (err error) {
	_emqx := h.Group("emqx")
	{
		_emqx.POST("auth", x.Emqx.Auth)
		_emqx.POST("acl", x.Emqx.Acl)
		_emqx.POST("bridge", x.Emqx.Bridge)
	}
	return
}

func (x *API) Initialize(ctx context.Context) (h *server.Hertz, err error) {
	h = x.Hertz

	if err = x.Transfer.Set(ctx, transfer.StreamOption{
		Key: "logset_imessages",
	}); err != nil {
		return
	}

	return
}
