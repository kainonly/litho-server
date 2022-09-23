package test

import (
	"context"
	"github.com/weplanx/server/api"
	"github.com/weplanx/server/bootstrap"
	"time"
)

func Initialize() (api *api.API, err error) {
	path := "./config/config.test.yml"
	values, err := bootstrap.LoadStaticValues(path)
	if err != nil {
		panic(err)
	}
	if api, err = bootstrap.NewAPI(values); err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err = api.Initialize(ctx); err != nil {
		return
	}

	return
}
