package model

import "time"

// Resource 资源
type Resource struct {
	ID        string    `gorm:"primaryKey;column:id;type:bigint"`
	CreatedAt time.Time `gorm:"column:created_at;not null"`
	UpdatedAt time.Time `gorm:"column:updated_at;not null"`
	Sort      int16     `gorm:"column:sort;not null;default:0"`             // 排序
	Active    bool      `gorm:"column:active;not null;default:true"`        // 状态
	Name      string    `gorm:"column:name;type:text;not null"`             // 资源名称
	Code      string    `gorm:"column:code;type:text;not null;uniqueIndex"` // 资源码
}

func (Resource) TableName() string {
	return "resource"
}
