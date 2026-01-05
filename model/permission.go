package model

import "time"

// Permission 特定授权表
type Permission struct {
	ID          string    `gorm:"primaryKey;column:id;type:bigint"`
	CreatedAt   time.Time `gorm:"column:created_at;type:timestamptz;not null"`
	UpdatedAt   time.Time `gorm:"column:updated_at;type:timestamptz;not null"`
	Active      bool      `gorm:"column:active;not null;default:true;comment:状态"`          // 状态
	Code        string    `gorm:"column:code;type:text;not null;uniqueIndex;comment:授权编码"` // 授权编码
	Description string    `gorm:"column:description;type:text;not null;comment:描述"`        // 描述
}

func (Permission) TableName() string {
	return "permission"
}
