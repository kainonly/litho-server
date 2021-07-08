package model

type RoleResourceRel struct {
	ID          uint64 `json:"-"`
	RoleKey     string `json:"-"`
	ResourceKey string `json:"-"`
}
