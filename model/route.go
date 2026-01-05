package model

import "time"

// Route 路由表
type Route struct {
	ID        string    `gorm:"primaryKey;column:id;type:bigint"`
	MenuID    string    `gorm:"column:menu_id;type:bigint;not null;index;comment:导航ID"` // 导航ID
	CreatedAt time.Time `gorm:"column:created_at;type:timestamptz;not null"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamptz;not null"`
	Sort      int16     `gorm:"column:sort;type:smallint;not null;default:0;comment:排序"`   // 排序
	Active    bool      `gorm:"column:active;not null;default:true;comment:状态"`            // 状态
	PID       string    `gorm:"column:pid;type:bigint;not null;default:0;comment:父级ID"`    // 父级ID
	Name      string    `gorm:"column:name;type:text;not null;comment:路由名称"`               // 路由名称
	Type      int16     `gorm:"column:type;type:smallint;not null;default:1;comment:路由类型"` // 路由类型: 1-页面, 2-分组, 3-外链, 4-iframe
	Icon      string    `gorm:"column:icon;type:text;not null;comment:字体图标"`               // 字体图标
	Link      string    `gorm:"column:link;type:text;not null;comment:链接"`                 // 链接
}

func (Route) TableName() string {
	return "route"
}
