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
			Path:   "dashboard",
			Name:   "仪表盘",
			Nav:    True(),
			Router: False(),
			Icon:   "dashboard",
		},
		{
			Parent: "dashboard",
			Path:   "dashboard/analysis",
			Name:   "分析页",
			Nav:    True(),
			Router: True(),
		},
		{
			Parent: "dashboard",
			Path:   "dashboard/monitor",
			Name:   "监控页",
			Nav:    True(),
			Router: True(),
		},
		{
			Parent: "dashboard",
			Path:   "dashboard/workbench",
			Name:   "工作台",
			Nav:    True(),
			Router: True(),
		},
		{
			Path:   "form",
			Name:   "表单页",
			Nav:    True(),
			Router: False(),
			Icon:   "form",
		},
		{
			Parent: "form",
			Path:   "form/basic",
			Name:   "基础表单",
			Nav:    True(),
			Router: True(),
		},
		{
			Parent: "form",
			Path:   "form/step",
			Name:   "分步表单",
			Nav:    True(),
			Router: True(),
		},
		{
			Parent: "form",
			Path:   "form/advanced",
			Name:   "高级表单",
			Nav:    True(),
			Router: True(),
		},
		{
			Path:   "list",
			Name:   "列表页",
			Nav:    True(),
			Router: False(),
			Icon:   "table",
		},
		{
			Parent: "list",
			Path:   "list/table",
			Name:   "查询表格",
			Nav:    True(),
			Router: True(),
		},
		{
			Parent: "list",
			Path:   "list/basic",
			Name:   "标准列表",
			Nav:    True(),
			Router: True(),
		},
		{
			Parent: "list",
			Path:   "list/card",
			Name:   "卡片列表",
			Nav:    True(),
			Router: True(),
		},
		{
			Path:   "profile",
			Name:   "详情页",
			Nav:    True(),
			Router: False(),
			Icon:   "profile",
		},
		{
			Parent: "profile",
			Path:   "profile/basic",
			Name:   "基础详情页",
			Nav:    True(),
			Router: True(),
		},
		{
			Parent: "profile",
			Path:   "profile/advanced",
			Name:   "高级详情页",
			Nav:    True(),
			Router: True(),
		},
		{
			Path:   "result",
			Name:   "结果页",
			Nav:    True(),
			Router: False(),
			Icon:   "check-circle",
		},
		{
			Parent: "result",
			Path:   "result/success",
			Name:   "成功页",
			Nav:    True(),
			Router: True(),
		},
		{
			Parent: "result",
			Path:   "result/fail",
			Name:   "失败页",
			Nav:    True(),
			Router: True(),
		},
		{
			Path:   "exception",
			Name:   "异常页",
			Nav:    True(),
			Router: False(),
			Icon:   "warning",
		},
		{
			Parent: "exception",
			Path:   "exception/403",
			Name:   "403",
			Nav:    True(),
			Router: True(),
		},
		{
			Parent: "exception",
			Path:   "exception/404",
			Name:   "404",
			Nav:    True(),
			Router: True(),
		},
		{
			Parent: "exception",
			Path:   "exception/500",
			Name:   "500",
			Nav:    True(),
			Router: True(),
		},
	}
	if err := db.Create(&data).Error; err != nil {
		t.Error(err)
	}
}
