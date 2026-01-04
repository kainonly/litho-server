package model

import "time"

// Role 权限
type Role struct {
	ID          string    `gorm:"primaryKey;column:id;type:bigint"`
	OrgID       string    `gorm:"column:org_id;type:bigint;not null"` // 所属组织ID
	CreatedAt   time.Time `gorm:"column:created_at;not null"`
	UpdatedAt   time.Time `gorm:"column:updated_at;not null"`
	Sort        int16     `gorm:"column:sort;not null;default:0"`        // 排序
	Active      bool      `gorm:"column:active;not null;default:true"`   // 状态
	Name        string    `gorm:"column:name;type:text;not null"`        // 权限名称
	Description string    `gorm:"column:description;type:text;not null"` // 权限描述
}

func (Role) TableName() string {
	return "role"
}
