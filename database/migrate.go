package main

import (
	"go.uber.org/fx"
	"lab-api/bootstrap"
	"lab-api/database/acl"
	"lab-api/database/admin"
	"lab-api/database/admin_basic"
	"lab-api/database/admin_role_rel"
	"lab-api/database/policy"
	"lab-api/database/resource"
	"lab-api/database/role"
	"lab-api/database/role_basic"
	"lab-api/database/role_policy"
	"lab-api/database/role_resource_rel"
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
