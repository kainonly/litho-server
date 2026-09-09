package main

import (
	"context"
	"fmt"
	"litho-api/api"
	"litho-api/bootstrap"
	"litho-api/common"
	"os"

	"github.com/cloudwego/hertz/pkg/app/server"
)

func main() {
	path := os.Getenv("CONFIG_PATH")
	if path == "" {
		path = "./config/values.yml"
	}

	if err := listen(path); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func listen(path string) (err error) {
	ctx := context.Background()
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
	h.Spin()
	return
}
