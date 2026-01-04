package model

import "time"

// Route 路由
type Route struct {
	ID        string    `gorm:"primaryKey;column:id;type:bigint"`
	MenuID    string    `gorm:"column:menu_id;type:bigint;not null"` // 导航ID
	CreatedAt time.Time `gorm:"column:created_at;not null"`
	UpdatedAt time.Time `gorm:"column:updated_at;not null"`
	Sort      int16     `gorm:"column:sort;not null;default:0"`            // 排序
	Active    bool      `gorm:"column:active;not null;default:true"`       // 状态
	PID       string    `gorm:"column:pid;type:bigint;not null;default:0"` // 父级ID
	Name      string    `gorm:"column:name;type:text;not null"`            // 路由名称
	Type      int16     `gorm:"column:type;not null;default:1"`            // 路由类型(1=页面,2=分组,3=外链,4=iframe)
	Icon      string    `gorm:"column:icon;type:text;not null"`            // 字体图标
	Link      string    `gorm:"column:link;type:text;not null"`            // 链接
}

func (Route) TableName() string {
	return "route"
}
