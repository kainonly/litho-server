package model

// RoleRoute 权限路由表
type RoleRoute struct {
	RoleID  string `gorm:"primaryKey;column:role_id;type:bigint"`  // 权限ID
	RouteID string `gorm:"primaryKey;column:route_id;type:bigint"` // 路由ID
}

func (RoleRoute) TableName() string {
	return "role_route"
}
