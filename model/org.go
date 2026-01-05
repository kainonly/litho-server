package model

import "time"

// Org 组织表
type Org struct {
	ID        string    `gorm:"primaryKey;column:id;type:bigint"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamptz;not null;index"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamptz;not null"`
	Active    bool      `gorm:"column:active;not null;default:true;comment:状态"` // 状态
	Name      string    `gorm:"column:name;type:text;not null;comment:组织名称"`    // 组织名称
}

func (Org) TableName() string {
	return "org"
}
