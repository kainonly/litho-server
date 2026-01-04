package model

// RouteResourceAction 路由资源操作
type RouteResourceAction struct {
	RouteID    string `gorm:"primaryKey;column:route_id;type:bigint"`    // 路由ID
	ResourceID string `gorm:"primaryKey;column:resource_id;type:bigint"` // 资源ID
	ActionID   string `gorm:"column:action_id;type:bigint"`              // 操作ID
}

func (RouteResourceAction) TableName() string {
	return "route_resource_action"
}
