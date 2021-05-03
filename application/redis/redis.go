package redis

import (
	"lab-api/application/redis/schema"
	"time"
)

type Model struct {
	Acl          *schema.Acl
	Resource     *schema.Resource
	Role         *schema.Role
	Admin        *schema.Admin
	RefreshToken *schema.RefreshToken
	Lock         *schema.Lock
}

func Initialize(dep schema.Dependency) *Model {
	c := new(Model)
	c.Acl = schema.NewAcl(dep)
	c.Resource = schema.NewResource(dep)
	c.Role = schema.NewRole(dep)
	c.Admin = schema.NewAdmin(dep)
	c.RefreshToken = schema.NewRefreshToken(dep)
	c.Lock = schema.NewLock(dep, 5, time.Second*900)
	return c
}
