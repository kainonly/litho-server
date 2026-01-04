package model

import "time"

// ResourceAction 资源操作
type ResourceAction struct {
	ID         string    `gorm:"primaryKey;column:id;type:bigint"`
	ResourceID string    `gorm:"column:resource_id;type:bigint;not null"` // 资源ID
	CreatedAt  time.Time `gorm:"column:created_at;not null"`
	UpdatedAt  time.Time `gorm:"column:updated_at;not null"`
	Sort       int16     `gorm:"column:sort;not null;default:0"`             // 排序
	Active     bool      `gorm:"column:active;not null;default:true"`        // 状态
	Name       string    `gorm:"column:name;type:text;not null"`             // 操作名称
	Code       string    `gorm:"column:code;type:text;not null;uniqueIndex"` // 操作码
}

func (ResourceAction) TableName() string {
	return "resource_action"
}
