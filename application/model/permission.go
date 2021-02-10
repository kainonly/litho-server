package model

type Permission struct {
	ID         uint64
	Key        string
	Name       string
	Note       string
	Status     bool
	CreateTime uint64 `gorm:"autoCreateTime"`
	UpdateTime uint64 `gorm:"autoUpdateTime"`
}
