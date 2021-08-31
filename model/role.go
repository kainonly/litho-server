package model

import "time"

type Role struct {
	ID          uint64     `json:"id"`
	Status      *bool      `gorm:"default:true" json:"status"`
	CreateTime  time.Time  `gorm:"autoCreateTime" json:"create_time"`
	UpdateTime  time.Time  `gorm:"autoUpdateTime" json:"update_time"`
	Code        string     `gorm:"type:varchar(20);not null;unique;comment:唯一码" json:"code"`
	Name        string     `gorm:"type:varchar(20);not null;comment:名称" json:"name"`
	Description string     `gorm:"type:text;comment:描述" json:"description"`
	Resources   []Resource `gorm:"many2many:role_resource_rel;foreignKey:Code;references:Path;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
}
