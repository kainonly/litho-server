package model

// RolePermission 权限特定授权表
type RolePermission struct {
	ID           string `gorm:"primaryKey;column:id;type:bigint"`
	RoleID       string `gorm:"column:role_id;type:bigint;not null;comment:权限ID"`         // 权限ID
	PermissionID string `gorm:"column:permission_id;type:bigint;not null;comment:特定授权ID"` // 特定授权ID
}

func (RolePermission) TableName() string {
	return "role_permission"
}
