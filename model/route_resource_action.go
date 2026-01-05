package model

// RouteResourceAction 路由资源表
type RouteResourceAction struct {
	ID         string `gorm:"primaryKey;column:id;type:bigint"`
	RouteID    string `gorm:"column:route_id;type:bigint;not null;uniqueIndex:idx_route_resource_action,priority:1;comment:路由ID"`    // 路由ID
	ResourceID string `gorm:"column:resource_id;type:bigint;not null;uniqueIndex:idx_route_resource_action,priority:2;comment:资源ID"` // 资源ID
	ActionID   string `gorm:"column:action_id;type:bigint;not null;uniqueIndex:idx_route_resource_action,priority:3;comment:操作ID"`   // 操作ID
}

func (RouteResourceAction) TableName() string {
	return "route_resource_action"
}
