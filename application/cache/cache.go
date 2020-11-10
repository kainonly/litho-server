package cache

import (
	"taste-api/application/cache/schema"
)

type Cache struct {
	Acl          *schema.Acl
	Resource     *schema.Resource
	Role         *schema.Role
	Admin        *schema.Admin
	RefreshToken *schema.RefreshToken
}

func Initialize(dep schema.Dependency) *Cache {
	c := new(Cache)
	c.Acl = schema.NewAcl(dep)
	c.Resource = schema.NewResource(dep)
	c.Role = schema.NewRole(dep)
	c.Admin = schema.NewAdmin(dep)
	c.RefreshToken = schema.NewRefreshToken(dep)
	return c
}
