package model

import "time"

// Menu 导航表
type Menu struct {
	ID        string    `gorm:"primaryKey;column:id;type:bigint"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamptz;not null;index"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamptz;not null"`
	Sort      int16     `gorm:"column:sort;type:smallint;not null;default:0;comment:排序"` // 排序
	Active    bool      `gorm:"column:active;not null;default:true;comment:状态"`          // 状态
	Name      string    `gorm:"column:name;type:text;not null;comment:导航名称"`             // 导航名称
}

func (Menu) TableName() string {
	return "menu"
}
