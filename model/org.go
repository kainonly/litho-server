package model

import "time"

// Org 组织
type Org struct {
	ID        string    `gorm:"primaryKey;column:id;type:bigint"`
	CreatedAt time.Time `gorm:"column:created_at;not null"`
	UpdatedAt time.Time `gorm:"column:updated_at;not null"`
	Active    bool      `gorm:"column:active;not null;default:true"` // 状态
	Name      string    `gorm:"column:name;type:text;not null"`      // 组织名称
}

func (Org) TableName() string {
	return "org"
}
