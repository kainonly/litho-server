package e2e

import (
	"context"
	"github.com/weplanx/api/api"
	"github.com/weplanx/api/bootstrap"
	"time"
)

func Initialize() (api *api.API, err error) {
	if api, err = bootstrap.NewAPI(); err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err = api.API.Initialize(ctx); err != nil {
		return
	}

	return
}
