package model

type AdminRoleRel struct {
	ID      uint64 `json:"-"`
	AdminId uint64 `json:"-"`
	RoleKey string `json:"-"`
}
