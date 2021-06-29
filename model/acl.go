package model

type Acl struct {
	ID         uint64
	Key        string `gorm:"type:varchar(200);unique"`
	Name       Object `gorm:"type:json"`
	Read       string `gorm:"type:longtext"`
	Write      string `gorm:"type:longtext"`
	Status     bool   `gorm:"default:1"`
	CreateTime uint64 `gorm:"autoCreateTime"`
	UpdateTime uint64 `gorm:"autoUpdateTime"`
}
