package model

import "time"

// Menu 导航菜单
type Menu struct {
	ID        string    `gorm:"primaryKey;column:id;type:bigint"`
	CreatedAt time.Time `gorm:"column:created_at;not null"`
	UpdatedAt time.Time `gorm:"column:updated_at;not null"`
	Sort      int16     `gorm:"column:sort;not null;default:0"`      // 排序
	Active    bool      `gorm:"column:active;not null;default:true"` // 状态
	Name      string    `gorm:"column:name;type:text;not null"`      // 导航名称
}

func (Menu) TableName() string {
	return "menu"
}
