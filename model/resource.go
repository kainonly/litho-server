package model

import "time"

// Resource 资源表
type Resource struct {
	ID        string    `gorm:"primaryKey;column:id;type:bigint"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamptz;not null;index"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamptz;not null"`
	Active    bool      `gorm:"column:active;not null;default:true;comment:状态"`        // 状态
	Name      string    `gorm:"column:name;type:text;not null;comment:资源名称"`           // 资源名称
	Code      string    `gorm:"column:code;type:text;not null;uniqueIndex;comment:路径"` // 路径
}

func (Resource) TableName() string {
	return "resource"
}
