package model

import (
	"server/common"
	"time"
)

type Route struct {
	ID         string     `gorm:"column:id;type:bigint;primaryKey" json:"id"`
	Status     *bool      `gorm:"column:status;type:boolean;not null;default:true" json:"status"`
	CreateTime *time.Time `gorm:"column:create_time;type:timestamptz;not null;default:now();index:,sort:desc" json:"create_time"`
	UpdateTime *time.Time `gorm:"column:update_time;type:timestamptz;not null;default:now()" json:"update_time"`
	Nav        string     `gorm:"column:nav;type:character varying;not null;index" json:"nav"`
	Type       int16      `gorm:"column:type;type:smallint;not null;index" json:"type"`
	Pid        string     `gorm:"column:pid;type:bigint;not null;index" json:"pid"`
	Name       string     `gorm:"column:name;type:character varying;not null" json:"name"`
	Icon       string     `gorm:"column:icon;type:character varying;not null" json:"icon"`
	Link       string     `gorm:"column:link;type:character varying;not null" json:"link"`
	Strategy   common.M   `gorm:"column:strategy;type:jsonb;not null;default:'{}';index:,type:gin" json:"strategy"`
	Sort       int16      `gorm:"column:sort;type:smallint;not null" json:"sort"`
}

func (Route) TableName() string {
	return "route"
}
