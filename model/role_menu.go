package model

// RoleMenu 权限导航表
type RoleMenu struct {
	RoleID string `gorm:"primaryKey;column:role_id;type:bigint"` // 权限ID
	MenuID string `gorm:"primaryKey;column:menu_id;type:bigint"` // 导航ID
}

func (RoleMenu) TableName() string {
	return "role_menu"
}
