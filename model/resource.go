package model

import "time"

type Resource struct {
	ID         uint64    `json:"id"`
	Status     *bool     `gorm:"default:true" json:"status"`
	CreateTime time.Time `gorm:"autoCreateTime" json:"create_time"`
	UpdateTime time.Time `gorm:"autoUpdateTime" json:"update_time"`
	Parent     uint64    `json:"parent"`
	Fragment   string    `gorm:"type:varchar(50);not null"`
	Name       string    `gorm:"type:varchar(20);not null"`
	Nav        *bool     `gorm:"default:false"`
	Router     *bool     `gorm:"default:false"`
	Icon       string    `gorm:"type:varchar(200)"`
}
