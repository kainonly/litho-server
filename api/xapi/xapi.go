package xapi

import (
	"go.uber.org/fx"
	"laboratory/api/xapi/admin"
	"laboratory/api/xapi/devops"
	"laboratory/api/xapi/page"
	"laboratory/api/xapi/role"
	"laboratory/api/xapi/schema"
	"laboratory/api/xapi/system"
)

var Provides = fx.Options(
	system.Provides,
	page.Provides,
	role.Provides,
	admin.Provides,
	devops.Provides,
	schema.Provides,
)
