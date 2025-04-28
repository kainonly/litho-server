package model

import (
	"time"
)

type Category struct {
	ID         string     `gorm:"column:id;type:bigint;primaryKey" json:"id"`
	CreateTime *time.Time `gorm:"column:create_time;type:timestamptz;not null;default:now();index:,sort:desc" json:"create_time"`
	UpdateTime *time.Time `gorm:"column:update_time;type:timestamptz;not null;default:now()" json:"update_time"`
	Type       int16      `gorm:"column:type;type:smallint;not null;index" json:"type"`
	Name       string     `gorm:"column:name;type:character varying;not null" json:"name"`
	Sort       int16      `gorm:"column:sort;type:smallint;not null" json:"sort"`
}

func (Category) TableName() string {
	return "category"
}
