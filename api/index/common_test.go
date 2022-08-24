package index_test

import (
	"github.com/weplanx/api/api"
	"github.com/weplanx/api/e2e"
	"os"
	"testing"
)

var x *api.API

func TestMain(m *testing.M) {
	os.Chdir("../../")
	x, _ = e2e.Initialize()
	os.Exit(m.Run())
}
