package main

import (
	"context"
	"github.com/weplanx/server/api"
	"github.com/weplanx/server/bootstrap"
	"time"
)

var x *api.API

func main() {
	var err error
	values, err := bootstrap.LoadStaticValues()
	if err != nil {
		panic(err)
	}
	if x, err = bootstrap.NewAPI(values); err != nil {
		panic(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	if err := setup(ctx); err != nil {
		panic(err)
	}
}

func setup(ctx context.Context) (err error) {
	if err = x.Values.Service.Reset(); err != nil {
		return
	}
	return
}
