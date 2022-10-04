package collections_test

import (
	"github.com/weplanx/server/api"
	"github.com/weplanx/server/bootstrap"
	"os"
	"testing"
)

var x *api.API

func TestMain(m *testing.M) {
	os.Chdir("../../")
	x, _ = bootstrap.UseTest()
	os.Exit(m.Run())
}
