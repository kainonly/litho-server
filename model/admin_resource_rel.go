package model

type AdminResourceRel struct {
	ID          uint64 `json:"-"`
	AdminId     uint64 `json:"-"`
	ResourceKey string `json:"-"`
}
