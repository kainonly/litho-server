package observability_test

import (
	"context"
	"github.com/weplanx/server/api"
	"github.com/weplanx/server/bootstrap"
	"github.com/weplanx/server/common"
	"os"
	"testing"
	"time"
)

var x *api.API

func TestMain(m *testing.M) {
	os.Chdir("../../")
	var err error
	var values *common.Values
	if values, err = bootstrap.LoadStaticValues(); err != nil {
		panic(err)
	}
	x, err = bootstrap.NewAPI(values)
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	if _, err = x.Initialize(ctx); err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}
