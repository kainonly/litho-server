package model

import (
	"taste-api/application/common/types"
)

type RoleBasic struct {
	ID         uint64
	Key        string     `gorm:"size:200;unique;not null;comment:role key"`
	Name       types.JSON `gorm:"type:json;not null;comment:role name"`
	Note       string     `gorm:"type:text;comment:note"`
	Status     bool       `gorm:"type:tinyint(1) unsigned;not null;default:1"`
	CreateTime uint64     `gorm:"not null;autoCreateTime"`
	UpdateTime uint64     `gorm:"not null;autoUpdateTime"`
}
