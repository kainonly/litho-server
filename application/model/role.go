package model

import (
	"github.com/kainonly/gin-extra/datatype"
)

type Role struct {
	ID         uint64
	Key        string
	Name       datatype.JSONObject `gorm:"type:jsonb"`
	Resource   datatype.JSONArray  `gorm:"type:json"`
	Acl        datatype.JSONArray  `gorm:"type:json"`
	Note       string
	Status     bool
	CreateTime uint64
	UpdateTime uint64
}
