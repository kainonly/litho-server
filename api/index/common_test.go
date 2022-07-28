package index_test

import (
	"context"
	"github.com/weplanx/server/api"
	"github.com/weplanx/server/bootstrap"
	"os"
	"testing"
	"time"
)

var x *api.API

func TestMain(m *testing.M) {
	os.Chdir("../../")
	var err error
	if x, err = bootstrap.NewAPI(); err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err = x.Initialize(ctx); err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}
