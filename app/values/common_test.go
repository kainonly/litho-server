package values

import (
	"api/common"
	"api/test"
	"os"
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
