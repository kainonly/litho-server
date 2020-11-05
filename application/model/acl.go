package model

import "taste-api/application/common/types"

type Acl struct {
	ID         uint64
	Keyid      string     `gorm:"size:200;unique;not null;comment:api access control key"`
	Name       types.JSON `gorm:"type:json;not null;comment:access control name"`
	Read       string     `gorm:"comment:list of readable api"`
	Write      string     `gorm:"comment:list of writable api"`
	Status     bool       `gorm:"type:tinyint(1) unsigned;not null;default:1"`
	CreateTime uint64     `gorm:"not null;autoCreateTime"`
	UpdateTime uint64     `gorm:"not null;autoUpdateTime"`
}
