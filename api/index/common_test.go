package index_test

import (
	"github.com/weplanx/server/api"
	"github.com/weplanx/server/test"
	"os"
	"testing"
)

var x *api.API

func TestMain(m *testing.M) {
	os.Chdir("../../")
	x, _ = test.Initialize()
	os.Exit(m.Run())
}
