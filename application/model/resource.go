package model

import "taste-api/application/common/datatype"

type Resource struct {
	ID         uint64
	Key        string              `gorm:"size:200;unique;not null;comment:resource access control key"`
	Parent     string              `gorm:"size:200;not null;default:origin;comment:parent node"`
	Name       datatype.JSONObject `gorm:"type:jsonb;not null;comment:access control name"`
	Nav        bool                `gorm:"not null;default:false;comment:show as navigation"`
	Router     bool                `gorm:"not null;default:false;comment:set as front-end routing"`
	Policy     bool                `gorm:"not null;default:false;comment:strategy node"`
	Icon       string              `gorm:"size:200;comment:iconfont"`
	Sort       uint8               `gorm:"comment:sort"`
	Status     bool                `gorm:"not null;default:true"`
	CreateTime uint64              `gorm:"not null;autoCreateTime"`
	UpdateTime uint64              `gorm:"not null;autoUpdateTime"`
}
