package model

import "time"

type Role struct {
	ID         uint64     `json:"id"`
	Status     *bool      `gorm:"default:true" json:"status"`
	CreateTime time.Time  `gorm:"autoCreateTime" json:"create_time"`
	UpdateTime time.Time  `gorm:"autoUpdateTime" json:"update_time"`
	Key        string     `gorm:""`
	Resource   []Resource `gorm:"many2many:role_resource_rel;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
