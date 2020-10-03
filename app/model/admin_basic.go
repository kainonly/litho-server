package model

type AdminBasic struct {
	ID         uint64
	Username   string `gorm:"size:30;not null;unique"`
	Password   string `gorm:"type:text;not null;"`
	Email      string `gorm:"size:200"`
	Phone      string `gorm:"size:20"`
	Call       string `gorm:"size:20"`
	Avatar     string `gorm:"type:text"`
	Status     bool   `gorm:"type:tinyint(1) unsigned;not null;default:1"`
	CreateTime uint64 `gorm:"not null;default:0;autoCreateTime"`
	UpdateTime uint64 `gorm:"not null;default:0;autoUpdateTime"`
}
