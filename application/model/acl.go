package model

import (
	"lab-api/application/common/datatype"
)

type Acl struct {
	ID         uint64
	Key        string              `gorm:"size:200;unique;not null;comment:api access control key"`
	Name       datatype.JSONObject `gorm:"type:jsonb;not null;comment:access control name"`
	Read       datatype.JSONArray  `gorm:"type:json;not null;default:'[]';comment:list of readable api"`
	Write      datatype.JSONArray  `gorm:"type:json;not null;default:'[]';comment:list of writable api"`
	Status     bool                `gorm:"not null;default:true"`
	CreateTime uint64              `gorm:"not null;autoCreateTime"`
	UpdateTime uint64              `gorm:"not null;autoUpdateTime"`
}
