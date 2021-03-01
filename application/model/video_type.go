package model

type VideoType struct {
	ID         uint64
	Name       string
	Sort       uint8
	Status     bool
	CreateTime uint64 `gorm:"autoCreateTime"`
	UpdateTime uint64 `gorm:"autoUpdateTime"`
}
