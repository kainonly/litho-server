package model

import "time"

type Admin struct {
	ID         uint64    `json:"id"`
	Username   string    `gorm:"type:varchar(20);not null;unique" json:"username"`
	Password   string    `gorm:"type:varchar(255);not null" json:"password"`
	Super      *bool     `gorm:"default:false"`
	Name       string    `gorm:"type:varchar(20)" json:"name"`
	Email      string    `gorm:"type:varchar(255)" json:"email"`
	Phone      string    `gorm:"type:varchar(20)" json:"phone"`
	Avatar     string    `gorm:"type:varchar(255)" json:"avatar"`
	Status     *bool     `gorm:"default:true" json:"status"`
	CreateTime time.Time `gorm:"autoCreateTime" json:"create_time"`
	UpdateTime time.Time `gorm:"autoUpdateTime" json:"update_time"`
}
