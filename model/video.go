package model

import (
	"time"
)

type Video struct {
	ID         string     `gorm:"column:id;type:bigint;primaryKey" json:"id"`
	CreateTime *time.Time `gorm:"column:create_time;type:timestamptz;not null;default:now();index:,sort:desc" json:"create_time"`
	UpdateTime *time.Time `gorm:"column:update_time;type:timestamptz;not null;default:now()" json:"update_time"`
	Name       string     `gorm:"column:name;type:character varying;not null" json:"name"`
	URL        string     `gorm:"column:url;type:character varying;not null" json:"url"`
}

func (Video) TableName() string {
	return "video"
}
