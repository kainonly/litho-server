package model

import (
	"server/common"
	"time"
)

type Role struct {
	ID          string     `gorm:"column:id;type:bigint;primaryKey" json:"id"`
	Status      *bool      `gorm:"column:status;type:boolean;not null;default:true;comment:状态" json:"status"`
	CreateTime  *time.Time `gorm:"column:create_time;type:timestamptz;not null;default:now();index:,sort:desc" json:"create_time"`
	UpdateTime  *time.Time `gorm:"column:update_time;type:timestamptz;not null;default:now()" json:"update_time"`
	Name        string     `gorm:"column:name;type:character varying;not null" json:"name"`
	Description string     `gorm:"column:description;type:character varying;not null" json:"description"`
	Strategy    common.M   `gorm:"column:strategy;type:jsonb;not null;default:'{}';index:,type:gin" json:"strategy"`
	Sort        int16      `gorm:"column:sort;type:smallint;not null" json:"sort"`
}

func (Role) TableName() string {
	return "role"
}
