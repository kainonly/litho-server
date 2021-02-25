package model

type Picture struct {
	ID         uint64
	TypeId     uint64
	Name       string
	Url        string
	Status     bool
	CreateTime uint64 `gorm:"autoCreateTime"`
	UpdateTime uint64 `gorm:"autoUpdateTime"`
}
