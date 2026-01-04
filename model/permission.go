package model

import "time"

// Permission 特定权限
type Permission struct {
	ID          string    `gorm:"primaryKey;column:id;type:bigint"`
	CreatedAt   time.Time `gorm:"column:created_at;not null"`
	UpdatedAt   time.Time `gorm:"column:updated_at;not null"`
	Active      bool      `gorm:"column:active;not null;default:true"`        // 状态
	Code        string    `gorm:"column:code;type:text;not null;uniqueIndex"` // 权限识别
	Description string    `gorm:"column:description;type:text;not null"`      // 描述
}

func (Permission) TableName() string {
	return "permission"
}
