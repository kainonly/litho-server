package model

// RoleRoute 权限路由表
type RoleRoute struct {
	ID      string `gorm:"primaryKey;column:id;type:bigint"`
	RoleID  string `gorm:"column:role_id;type:bigint;not null;uniqueIndex:idx_role_route,priority:1;comment:权限ID"`  // 权限ID
	RouteID string `gorm:"column:route_id;type:bigint;not null;uniqueIndex:idx_role_route,priority:2;comment:路由ID"` // 路由ID
}

func (RoleRoute) TableName() string {
	return "role_route"
}
