package e2e

import (
	"context"
	"fmt"
	"github.com/weplanx/server/api"
	"github.com/weplanx/server/bootstrap"
	"time"
)

func Initialize() (api *api.API, err error) {
	if api, err = bootstrap.NewAPI(); err != nil {
		return
	}
	api.Values.Namespace = fmt.Sprintf("%s_test", api.Values.Namespace)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err = api.Initialize(ctx); err != nil {
		return
	}

	return
}
