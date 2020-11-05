package model

import "taste-api/application/common/types"

type Resource struct {
	ID         uint64
	Keyid      string     `gorm:"size:200;unique;not null;comment:resource access control key"`
	Parent     string     `gorm:"size:200;not null;default:origin;comment:parent node"`
	Name       types.JSON `gorm:"type:json;not null;comment:access control name"`
	Nav        bool       `gorm:"type:tinyint(1) unsigned;not null;default:0;comment:show as navigation"`
	Router     bool       `gorm:"type:tinyint(1) unsigned;not null;default:0;comment:set as front-end routing"`
	Policy     bool       `gorm:"type:tinyint(1) unsigned;not null;default:0;comment:strategy node"`
	Icon       string     `gorm:"size:200;comment:iconfont"`
	Sort       uint8      `gorm:"comment:sort"`
	Status     bool       `gorm:"type:tinyint(1) unsigned;not null;default:1"`
	CreateTime uint64     `gorm:"not null;autoCreateTime"`
	UpdateTime uint64     `gorm:"not null;autoUpdateTime"`
}
