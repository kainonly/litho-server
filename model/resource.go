package model

import (
	"database/sql/driver"
	jsoniter "github.com/json-iterator/go"
	"gorm.io/gorm"
)

type Resource struct {
	ID     int64  `json:"id"`
	Name   string `gorm:"type:varchar;not null" json:"name"`
	Path   string `gorm:"type:varchar;not null;unique" json:"path"`
	Parent string `gorm:"type:varchar;default:'root'" json:"parent"`
	Router Router `gorm:"type:jsonb;default:'{}'" json:"router"`
	Nav    *bool  `gorm:"default:false" json:"nav"`
	Icon   string `gorm:"type:varchar" json:"icon"`
	Sort   int8   `gorm:"default:0" json:"sort"`
}

type Router struct {
	Template string       `json:"template,omitempty"`
	Schema   string       `json:"schema,omitempty"`
	Option   RouterOption `json:"options,omitempty"`
}

type RouterOption struct {
	Fetch   bool         `json:"fetch,omitempty"`
	Columns []ViewColumn `json:"columns,omitempty"`
}

type ViewColumn struct {
	Name string `json:"name"`
}

func (x *Router) Scan(input interface{}) error {
	return jsoniter.Unmarshal(input.([]byte), x)
}

func (x Router) Value() (driver.Value, error) {
	return jsoniter.Marshal(x)
}

func GenerateResources(tx *gorm.DB) (err error) {
	if tx.Migrator().HasTable(&Resource{}) {
		if err = tx.Migrator().DropTable(&Resource{}); err != nil {
			return
		}
	}
	if err = tx.AutoMigrate(&Resource{}); err != nil {
		return
	}
	tx.Exec("create index router_gin on resource using gin(router)")
	data := []Resource{
		{
			Name:   "仪表盘",
			Path:   "dashboard",
			Parent: "root",
			Router: Router{
				Template: "manual",
			},
			Nav:  True(),
			Icon: "dashboard",
		},
		{
			Name:   "个人中心",
			Path:   "center",
			Parent: "root",
		},
		{
			Name:   "我的信息",
			Path:   "center/profile",
			Parent: "center",
			Router: Router{
				Template: "manual",
			},
		},
		{
			Name:   "消息通知",
			Path:   "center/notification",
			Parent: "center",
			Router: Router{
				Template: "manual",
			},
		},
		{
			Name:   "设置",
			Path:   "setting",
			Parent: "root",
			Icon:   "setting",
		},
		{
			Name:   "资源管理",
			Path:   "setting/resource",
			Parent: "setting",
			Router: Router{
				Template: "manual",
			},
			Nav: True(),
		},
		{
			Name:   "权限管理",
			Path:   "setting/role",
			Parent: "setting",
			Router: Router{
				Template: "list",
				Schema:   "role",
			},
			Nav: True(),
		},
		{
			Name:   "创建资源",
			Path:   "setting/role/create",
			Parent: "setting/role",
			Router: Router{
				Template: "page",
				Schema:   "role",
			},
		},
		{
			Name:   "更新资源",
			Path:   "setting/role/update",
			Parent: "setting/role",
			Router: Router{
				Template: "page",
				Schema:   "role",
			},
		},
		{
			Name:   "成员管理",
			Path:   "setting/admin",
			Parent: "setting",
			Router: Router{
				Template: "list",
				Schema:   "admin",
			},
			Nav: True(),
		},
		{
			Name:   "创建资源",
			Path:   "setting/admin/create",
			Parent: "setting/admin",
			Router: Router{
				Template: "page",
				Schema:   "admin",
			},
		},
		{
			Name:   "更新资源",
			Path:   "setting/admin/update",
			Parent: "setting/admin",
			Router: Router{
				Template: "page",
				Schema:   "admin",
			},
		},
	}
	return tx.Create(&data).Error
}
