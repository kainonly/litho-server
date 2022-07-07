package values

import (
	"os"
	"server/common"
	"server/test"
	"testing"
)

var i *common.Inject
var service *Service

func TestMain(m *testing.M) {
	os.Chdir("../../")
	i, _ = test.Bed()
	service = &Service{Inject: i}
	os.Exit(m.Run())
}
