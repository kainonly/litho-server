package main

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app/server"
	"os"
	"server/api"
	"server/bootstrap"
	"server/common"
)

func main() {
	if err := listen("./config/values.yml"); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func listen(path string) (err error) {
	ctx := context.TODO()
	var v *common.Values
	if v, err = bootstrap.LoadStaticValues(path); err != nil {
		return
	}
	var x *api.API
	if x, err = bootstrap.NewAPI(v); err != nil {
		return
	}

	var h *server.Hertz
	if h, err = x.Initialize(ctx); err != nil {
		return
	}
	if err = x.SetupRoutes(h); err != nil {
		return
	}
	h.Spin()
	return
}
