package model

import "time"

// Role 权限表
type Role struct {
	ID          string    `gorm:"primaryKey;column:id;type:bigint"`
	OrgID       string    `gorm:"column:org_id;type:bigint;not null;index;comment:所属组织ID"` // 所属组织ID
	CreatedAt   time.Time `gorm:"column:created_at;type:timestamptz;not null;index"`
	UpdatedAt   time.Time `gorm:"column:updated_at;type:timestamptz;not null"`
	Sort        int16     `gorm:"column:sort;type:smallint;not null;default:0;comment:排序"` // 排序
	Active      bool      `gorm:"column:active;not null;default:true;comment:状态"`          // 状态
	Name        string    `gorm:"column:name;type:text;not null;comment:权限名称"`             // 权限名称
	Description string    `gorm:"column:description;type:text;not null;comment:权限描述"`      // 权限描述
}

func (Role) TableName() string {
	return "role"
}
