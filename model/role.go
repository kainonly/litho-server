package model

import "time"

type Role struct {
	ID          uint64       `json:"id"`
	Status      *bool        `gorm:"default:true" json:"status"`
	CreateTime  time.Time    `gorm:"autoCreateTime" json:"create_time"`
	UpdateTime  time.Time    `gorm:"autoUpdateTime" json:"update_time"`
	Name        string       `gorm:"type:varchar(20);not null;comment:名称"`
	Description string       `gorm:"type:text;comment:描述"`
	Resources   []Resource   `gorm:"many2many:role_resource_rel;constraint:OnDelete:CASCADE"`
	Permission  []Permission `gorm:"many2many:role_permission_rel;constraint:OnDelete:CASCADE"`
}
