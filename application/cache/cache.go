package cache

import (
	"taste-api/application/cache/schema"
)

type Cache struct {
	*schema.Acl
	*schema.Resource
	*schema.Role
	*schema.Admin
}

func Initialize(dep schema.Dependency) *Cache {
	c := new(Cache)
	c.Acl = schema.NewAcl(dep)
	c.Resource = schema.NewResource(dep)
	c.Role = schema.NewRole(dep)
	c.Admin = schema.NewAdmin(dep)
	return c
}
