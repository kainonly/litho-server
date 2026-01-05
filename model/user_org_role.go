package model

// UserOrgRole 用户组织权限表
type UserOrgRole struct {
	ID     string `gorm:"primaryKey;column:id;type:bigint"`
	UserID string `gorm:"column:user_id;type:bigint;not null;uniqueIndex:idx_user_org_role;comment:用户ID"` // 用户ID
	OrgID  string `gorm:"column:org_id;type:bigint;not null;uniqueIndex:idx_user_org_role;comment:组织ID"`  // 组织ID
	RoleID string `gorm:"column:role_id;type:bigint;not null;index;comment:权限ID"`                         // 权限ID
}

func (UserOrgRole) TableName() string {
	return "user_org_role"
}
