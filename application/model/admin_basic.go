package model

type AdminBasic struct {
	ID         uint64
	Username   string `gorm:"size:50;not null;unique"`
	Password   string `gorm:"not null"`
	Email      string `gorm:"size:200"`
	Phone      string `gorm:"size:20"`
	Call       string `gorm:"size:20"`
	Avatar     string `gorm:"type:text"`
	Status     bool   `gorm:"not null;default:true"`
	CreateTime uint64 `gorm:"not null;autoCreateTime"`
	UpdateTime uint64 `gorm:"not null;autoUpdateTime"`
}
