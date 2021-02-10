package model

import "github.com/kainonly/gin-extra/datatype"

type Resource struct {
	ID         uint64
	Key        string
	Parent     string
	Name       datatype.JSONObject `gorm:"type:json"`
	Nav        bool
	Router     bool
	Policy     bool
	Icon       string
	Sort       uint8
	Status     bool
	CreateTime uint64 `gorm:"autoCreateTime"`
	UpdateTime uint64 `gorm:"autoUpdateTime"`
}
