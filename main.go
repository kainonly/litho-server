package main

import (
	"context"
	"github.com/weplanx/server/bootstrap"
	"time"
)

func main() {
	path := "./config/config.yml"
	values, err := bootstrap.LoadStaticValues(path)
	if err != nil {
		panic(err)
	}

	api, err := bootstrap.NewAPI(values)
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	h, err := api.Initialize(ctx)
	if err != nil {
		panic(err)
	}

	if err = api.Routes(h); err != nil {
		panic(err)
	}

	h.Spin()
}
