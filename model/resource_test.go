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
			ID:       1,
			Parent:   0,
			Fragment: "dashboard",
			Name:     "仪表盘",
			Nav:      True(),
			Router:   False(),
			Icon:     "dashboard",
		},
		{
			ID:       2,
			Parent:   1,
			Fragment: "analysis",
			Name:     "分析页",
			Nav:      True(),
			Router:   True(),
			Icon:     "",
		},
		{
			ID:       3,
			Parent:   1,
			Fragment: "monitor",
			Name:     "监控页",
			Nav:      True(),
			Router:   True(),
			Icon:     "",
		},
		{
			ID:       4,
			Parent:   1,
			Fragment: "workbench",
			Name:     "工作台",
			Nav:      True(),
			Router:   True(),
			Icon:     "",
		},
		{
			ID:       5,
			Parent:   0,
			Fragment: "form",
			Name:     "表单页",
			Nav:      True(),
			Router:   False(),
			Icon:     "form",
		},
		{
			ID:       6,
			Parent:   5,
			Fragment: "basic",
			Name:     "基础表单",
			Nav:      True(),
			Router:   True(),
			Icon:     "",
		},
		{
			ID:       7,
			Parent:   5,
			Fragment: "step",
			Name:     "分步表单",
			Nav:      True(),
			Router:   True(),
			Icon:     "",
		},
		{
			ID:       8,
			Parent:   5,
			Fragment: "advanced",
			Name:     "高级表单",
			Nav:      True(),
			Router:   True(),
			Icon:     "",
		},
		{
			ID:       9,
			Parent:   0,
			Fragment: "list",
			Name:     "列表页",
			Nav:      True(),
			Router:   False(),
			Icon:     "table",
		},
		{
			ID:       10,
			Parent:   9,
			Fragment: "table",
			Name:     "查询表格",
			Nav:      True(),
			Router:   True(),
			Icon:     "",
		},
		{
			ID:       11,
			Parent:   9,
			Fragment: "basic",
			Name:     "标准列表",
			Nav:      True(),
			Router:   True(),
			Icon:     "",
		},
		{
			ID:       12,
			Parent:   9,
			Fragment: "card",
			Name:     "卡片列表",
			Nav:      True(),
			Router:   True(),
			Icon:     "",
		},
		{
			ID:       13,
			Parent:   0,
			Fragment: "profile",
			Name:     "详情页",
			Nav:      True(),
			Router:   False(),
			Icon:     "profile",
		},
		{
			ID:       14,
			Parent:   13,
			Fragment: "basic",
			Name:     "基础详情页",
			Nav:      True(),
			Router:   True(),
			Icon:     "",
		},
		{
			ID:       15,
			Parent:   13,
			Fragment: "advanced",
			Name:     "高级详情页",
			Nav:      True(),
			Router:   True(),
			Icon:     "",
		},
		{
			ID:       16,
			Parent:   0,
			Fragment: "result",
			Name:     "结果页",
			Nav:      True(),
			Router:   False(),
			Icon:     "check-circle",
		},
		{
			ID:       17,
			Parent:   16,
			Fragment: "success",
			Name:     "成功页",
			Nav:      True(),
			Router:   True(),
			Icon:     "",
		},
		{
			ID:       18,
			Parent:   16,
			Fragment: "fail",
			Name:     "失败页",
			Nav:      True(),
			Router:   True(),
			Icon:     "",
		},
		{
			ID:       19,
			Parent:   0,
			Fragment: "exception",
			Name:     "异常页",
			Nav:      True(),
			Router:   False(),
			Icon:     "warning",
		},
		{
			ID:       20,
			Parent:   19,
			Fragment: "73",
			Name:     "73",
			Nav:      True(),
			Router:   True(),
			Icon:     "",
		},
		{
			ID:       21,
			Parent:   19,
			Fragment: "74",
			Name:     "74",
			Nav:      True(),
			Router:   True(),
			Icon:     "",
		},
		{
			ID:       22,
			Parent:   19,
			Fragment: "500",
			Name:     "500",
			Nav:      True(),
			Router:   True(),
			Icon:     "",
		},
	}
	if err := db.Create(&data).Error; err != nil {
		t.Error(err)
	}
	// 指定ID写入后手动更新序列值
	sql := `alter sequence resource_id_seq restart 23`
	if err := db.Exec(sql).Error; err != nil {
		t.Error(err)
	}
}
