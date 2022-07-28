package sessions_test

import (
	"github.com/weplanx/server/api"
	"github.com/weplanx/server/e2e"
	"os"
	"testing"
)

var x *api.API

func TestMain(m *testing.M) {
	os.Chdir("../../")
	x, _ = e2e.Initialize()
	os.Exit(m.Run())
}
