package main

import (
	"go.uber.org/fx"
	"taste-api/bootstrap"
	"taste-api/database/acl"
	"taste-api/database/admin"
	"taste-api/database/admin_basic"
	"taste-api/database/admin_role_rel"
	"taste-api/database/policy"
	"taste-api/database/resource"
	"taste-api/database/role"
	"taste-api/database/role_basic"
	"taste-api/database/role_policy"
	"taste-api/database/role_resource_rel"
)

func main() {
	fx.New(
		fx.Provide(
			bootstrap.LoadConfiguration,
			bootstrap.InitializeDatabase,
		),
		fx.Invoke(
			acl.Setup,
			resource.Setup,
			policy.Setup,
			role_basic.Setup,
			role_resource_rel.Setup,
			admin_basic.Setup,
			admin_role_rel.Setup,
			role_policy.Setup,
			role.Setup,
			admin.Setup,
		),
	).Done()
}
