package model

import "time"

type Role struct {
	ID         uint64     `json:"id"`
	Status     *bool      `gorm:"default:true" json:"status"`
	CreateTime time.Time  `gorm:"autoCreateTime" json:"create_time"`
	UpdateTime time.Time  `gorm:"autoUpdateTime" json:"update_time"`
	Key        string     `gorm:"type:varchar(20);not null"`
	Name       string     `gorm:"type:varchar(20);not null"`
	Note       string     `gorm:"type:text"`
	Resource   []Resource `gorm:"many2many:role_to_resource"`
}

type RoleToResource struct {
	ID         uint64
	RoleID     uint64 `gorm:"primaryKey"`
	ResourceID uint64 `gorm:"primaryKey"`
}
