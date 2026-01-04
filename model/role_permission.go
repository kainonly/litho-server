package model

// RolePermission 权限特定权限表
type RolePermission struct {
	RoleID       string `gorm:"primaryKey;column:role_id;type:bigint"`       // 权限ID
	PermissionID string `gorm:"primaryKey;column:permission_id;type:bigint"` // 特定权限ID
}

func (RolePermission) TableName() string {
	return "role_permission"
}
