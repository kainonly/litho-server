package model

import (
	"github.com/kainonly/gin-extra/datatype"
)

type RoleMix struct {
	ID         uint64
	Key        string
	Name       datatype.JSONObject `gorm:"type:json"`
	Resource   string
	Acl        string
	Permission string
	Note       string
	Status     bool
	CreateTime uint64
	UpdateTime uint64
}
