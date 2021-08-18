package service

type Role struct {
	*Dependency
	Key string
}

func NewRole(d Dependency) *Role {
	return &Role{
		Dependency: &d,
		Key:        d.Config.RedisKey("role"),
	}
}
