package model

// RoleMenu 权限导航表
type RoleMenu struct {
	ID     string `gorm:"primaryKey;column:id;type:bigint"`
	RoleID string `gorm:"column:role_id;type:bigint;not null;uniqueIndex:idx_role_menu,priority:1;comment:权限ID"` // 权限ID
	MenuID string `gorm:"column:menu_id;type:bigint;not null;uniqueIndex:idx_role_menu,priority:2;comment:导航ID"` // 导航ID
}

func (RoleMenu) TableName() string {
	return "role_menu"
}
