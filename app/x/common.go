package x

import (
	"api/app/x/admin"
	"api/app/x/devops"
	"api/app/x/page"
	"api/app/x/schema"
	"github.com/google/wire"
)

var Provides = wire.NewSet(
	wire.Struct(new(InjectController), "*"),
	wire.Struct(new(InjectService), "*"),
	NewController,
	NewService,
	page.Provides,
	admin.Provides,
	devops.Provides,
	schema.Provides,
)
