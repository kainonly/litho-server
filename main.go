package main

import (
	"context"
	"github.com/weplanx/server/bootstrap"
	"time"
)

func main() {
	api, err := bootstrap.NewAPI()
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err = api.Initialize(ctx); err != nil {
		panic(err)
	}

	h, err := api.Run()
	if err != nil {
		panic(err)
	}

	if _, err = api.Routes(h); err != nil {
		panic(err)
	}

	h.Spin()
}
