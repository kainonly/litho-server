package main

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/weplanx/server/api"
	"github.com/weplanx/server/bootstrap"
	"github.com/weplanx/server/common"
)

func main() {
	var err error
	var values *common.Values
	ctx := context.Background()
	values, err = bootstrap.LoadStaticValues("./config/values.yml")
	var x *api.API
	if x, err = bootstrap.NewAPI(values); err != nil {
		return
	}
	var h *server.Hertz
	if h, err = x.Initialize(ctx); err != nil {
		return
	}
	if err = x.Routes(h); err != nil {
		return
	}
	h.Spin()
}
