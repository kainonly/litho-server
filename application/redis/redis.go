package redis

import (
	"lab-api/application/redis/schema"
)

type Model struct {
	Acl          *schema.Acl
	Resource     *schema.Resource
	Role         *schema.Role
	Admin        *schema.Admin
	RefreshToken *schema.RefreshToken
	UserLock     *schema.UserLock
}

func Initialize(dep schema.Dependency) *Model {
	c := new(Model)
	c.Acl = schema.NewAcl(dep)
	c.Resource = schema.NewResource(dep)
	c.Role = schema.NewRole(dep)
	c.Admin = schema.NewAdmin(dep)
	c.RefreshToken = schema.NewRefreshToken(dep)
	c.UserLock = schema.NewUserLock(dep)
	return c
}
