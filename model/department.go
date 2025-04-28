package model

import (
	"time"
)

type Department struct {
	ID          string     `gorm:"column:id;type:bigint;primaryKey" json:"id"`
	Status      *bool      `gorm:"column:status;type:boolean;not null;default:true" json:"status"`
	CreateTime  *time.Time `gorm:"column:create_time;type:timestamptz;not null;default:now();index:,sort:desc" json:"create_time"`
	UpdateTime  *time.Time `gorm:"column:update_time;type:timestamptz;not null;default:now()" json:"update_time"`
	Name        string     `gorm:"column:name;type:character varying;not null;index" json:"name"`
	Description string     `gorm:"column:description;type:character varying;not null" json:"description"`
}

func (Department) TableName() string {
	return "department"
}
