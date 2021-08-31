package model

import (
	"testing"
)

func TestResource(t *testing.T) {
	if err := db.AutoMigrate(&Resource{}); err != nil {
		t.Error(err)
	}
	data := []Resource{
		{
			Path:   "workbench",
			Name:   "工作台",
			Nav:    True(),
			Router: False(),
			Icon:   "desktop",
		},
		{
			Parent: "workbench",
			Path:   "workbench/dashboard",
			Name:   "仪表盘",
			Nav:    True(),
			Router: True(),
		},
		{
			Parent: "workbench",
			Path:   "workbench/message",
			Name:   "消息中心",
			Nav:    True(),
			Router: True(),
		},
		{
			Parent: "workbench",
			Path:   "workbench/profile",
			Name:   "个人中心",
			Nav:    True(),
			Router: True(),
		},
		{
			Path:   "settings",
			Name:   "设置",
			Nav:    True(),
			Router: False(),
			Icon:   "setting",
		},
		{
			Parent: "settings",
			Path:   "settings/team",
			Name:   "组织成员",
			Nav:    True(),
			Router: True(),
		},
		{
			Parent: "settings",
			Path:   "settings/resource",
			Name:   "资源控制",
			Nav:    True(),
			Router: True(),
		},
	}
	if err := db.Create(&data).Error; err != nil {
		t.Error(err)
	}
}
