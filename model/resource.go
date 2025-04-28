package model

import (
	"server/common"
	"time"
)

type Resource struct {
	ID         string         `gorm:"column:id;type:bigint;primaryKey" json:"id"`
	CreateTime *time.Time     `gorm:"column:create_time;type:timestamptz;not null;default:now();index:,sort:desc" json:"create_time"`
	UpdateTime *time.Time     `gorm:"column:update_time;type:timestamptz;not null;default:now()" json:"update_time"`
	Name       string         `gorm:"column:name;type:character varying;not null" json:"name"`
	Path       string         `gorm:"column:path;type:character varying;not null;index" json:"path"`
	Actions    common.Actions `gorm:"column:actions;type:jsonb;not null;default:'[]'" json:"actions"`
}

func (Resource) TableName() string {
	return "resource"
}
