package model

import "time"

// ResourceAction 资源操作表
type ResourceAction struct {
	ID         string    `gorm:"primaryKey;column:id;type:bigint"`
	ResourceID string    `gorm:"column:resource_id;type:bigint;not null;index;comment:资源ID"` // 资源ID
	CreatedAt  time.Time `gorm:"column:created_at;type:timestamptz;not null;index"`
	UpdatedAt  time.Time `gorm:"column:updated_at;type:timestamptz;not null"`
	Active     bool      `gorm:"column:active;not null;default:true;comment:状态"`        // 状态
	Name       string    `gorm:"column:name;type:text;not null;comment:操作名称"`           // 操作名称
	Code       string    `gorm:"column:code;type:text;not null;uniqueIndex;comment:编码"` // 编码
}

func (ResourceAction) TableName() string {
	return "resource_action"
}
