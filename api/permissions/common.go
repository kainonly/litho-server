package permissions

import (
	"litho-api/common"

	"github.com/goforj/wire"
)

var Provides = wire.NewSet(
	wire.Struct(new(Controller), "*"),
	wire.Struct(new(Service), "*"),
)

const Key = "permissions"

type Controller struct {
	V *common.Values

	PermissionsX *Service
}

type Service struct {
	*common.Inject
}

type M = map[string]any
