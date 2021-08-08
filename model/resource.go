package model

type Resource struct {
	Common

	Parent   uint64   `json:"parent"`
	Fragment string   `gorm:"type:varchar(50);not null"`
	Acl      []string `gorm:"type:varchar(50)[]"`
	Name     string   `gorm:"type:varchar(20);not null"`
	Nav      *bool    `gorm:"default:false"`
	Router   *bool    `gorm:"default:false"`
	Icon     string   `gorm:"type:varchar(200)"`
}
