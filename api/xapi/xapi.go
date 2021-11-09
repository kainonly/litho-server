package xapi

import (
	"github.com/google/wire"
	"laboratory/api/xapi/admin"
	"laboratory/api/xapi/devops"
	"laboratory/api/xapi/page"
	"laboratory/api/xapi/schema"
	"laboratory/api/xapi/system"
)

var Provides = wire.NewSet(
	system.Provides,
	page.Provides,
	admin.Provides,
	devops.Provides,
	schema.Provides,
)
