package model

// UserOrgRole 用户组织权限表
type UserOrgRole struct {
	UserID string `gorm:"primaryKey;column:user_id;type:bigint"` // 用户ID
	OrgID  string `gorm:"primaryKey;column:org_id;type:bigint"`  // 组织ID
	RoleID string `gorm:"column:role_id;type:bigint"`            // 权限ID
}

func (UserOrgRole) TableName() string {
	return "user_org_role"
}
