package model

type Policy struct {
	ID          uint64 `json:"-"`
	Policy      *bool  `json:"-"`
	ResourceKey string `json:"-"`
	AclKey      string `json:"-"`
}
