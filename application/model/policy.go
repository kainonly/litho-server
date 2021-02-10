package model

type Policy struct {
	ID          uint64
	ResourceKey string
	AclKey      string
	Policy      uint8
}
