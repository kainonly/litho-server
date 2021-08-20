package model

import (
	"gorm.io/gorm"
	"testing"
)

func TestResource(t *testing.T) {
	if err := db.SetupJoinTable(&Resource{}, "ResourceAclRel", &ResourceAclRel{}); err != nil {
		t.Error(err)
	}
	if err := db.AutoMigrate(&Resource{}); err != nil {
		t.Error(err)
	}
	if err := db.Transaction(func(tx *gorm.DB) (err error) {
		system := Resource{
			Parent:   0,
			Fragment: "system",
			Name:     "系统设置",
			Nav:      True(),
			Icon:     "setting",
			Sort:     1,
		}
		if err = tx.Create(&system).Error; err != nil {
			return
		}
		systemItems := []Resource{
			{
				Parent:   system.ID,
				Fragment: "acl",
				Name:     "访问控制",
				Nav:      True(),
				Router:   True(),
				Strategy: True(),
			},
			{
				Parent:   system.ID,
				Fragment: "resource",
				Name:     "资源控制",
				Nav:      True(),
				Router:   True(),
				Strategy: True(),
			},
			{
				Parent:   system.ID,
				Fragment: "permission",
				Name:     "特殊授权",
				Nav:      True(),
				Router:   True(),
				Strategy: True(),
			},
			{
				Parent:   system.ID,
				Fragment: "role",
				Name:     "权限组",
				Nav:      True(),
				Router:   True(),
				Strategy: True(),
			},
			{
				Parent:   system.ID,
				Fragment: "admin",
				Name:     "成员管理",
				Nav:      True(),
				Router:   True(),
				Strategy: True(),
			},
		}
		if err = tx.Create(&systemItems).Error; err != nil {
			return
		}
		dashboard := Resource{
			Parent:   0,
			Fragment: "dashboard",
			Name:     "仪表盘",
			Nav:      True(),
			Router:   False(),
			Icon:     "dashboard",
		}
		if err = tx.Create(&dashboard).Error; err != nil {
			return
		}
		dashboardItems := []Resource{
			{
				Parent:   dashboard.ID,
				Fragment: "analysis",
				Name:     "分析页",
				Nav:      True(),
				Router:   True(),
			},
			{
				Parent:   dashboard.ID,
				Fragment: "monitor",
				Name:     "监控页",
				Nav:      True(),
				Router:   True(),
			},
			{
				Parent:   dashboard.ID,
				Fragment: "workbench",
				Name:     "工作台",
				Nav:      True(),
				Router:   True(),
			},
		}
		if err = tx.Create(&dashboardItems).Error; err != nil {
			return
		}
		form := Resource{
			Parent:   0,
			Fragment: "form",
			Name:     "表单页",
			Nav:      True(),
			Router:   False(),
			Icon:     "form",
		}
		if err = tx.Create(&form).Error; err != nil {
			return
		}
		formItems := []Resource{
			{
				Parent:   form.ID,
				Fragment: "basic",
				Name:     "基础表单",
				Nav:      True(),
				Router:   True(),
			},
			{
				Parent:   form.ID,
				Fragment: "step",
				Name:     "分步表单",
				Nav:      True(),
				Router:   True(),
			},
			{
				Parent:   form.ID,
				Fragment: "advanced",
				Name:     "高级表单",
				Nav:      True(),
				Router:   True(),
			},
		}
		if err = tx.Create(&formItems).Error; err != nil {
			return
		}
		list := Resource{
			Parent:   0,
			Fragment: "list",
			Name:     "列表页",
			Nav:      True(),
			Router:   False(),
			Icon:     "table",
		}
		if err = tx.Create(&list).Error; err != nil {
			return
		}
		listItems := []Resource{
			{
				Parent:   list.ID,
				Fragment: "table",
				Name:     "查询表格",
				Nav:      True(),
				Router:   True(),
			},
			{
				Parent:   list.ID,
				Fragment: "basic",
				Name:     "标准列表",
				Nav:      True(),
				Router:   True(),
			},
			{
				Parent:   list.ID,
				Fragment: "card",
				Name:     "卡片列表",
				Nav:      True(),
				Router:   True(),
			},
		}
		if err = tx.Create(&listItems).Error; err != nil {
			return
		}
		profile := Resource{
			Parent:   0,
			Fragment: "profile",
			Name:     "详情页",
			Nav:      True(),
			Router:   False(),
			Icon:     "profile",
		}
		if err = tx.Create(&profile).Error; err != nil {
			return
		}
		profileItems := []Resource{
			{
				Parent:   profile.ID,
				Fragment: "basic",
				Name:     "基础详情页",
				Nav:      True(),
				Router:   True(),
			},
			{
				Parent:   profile.ID,
				Fragment: "advanced",
				Name:     "高级详情页",
				Nav:      True(),
				Router:   True(),
			},
		}
		if err = tx.Create(&profileItems).Error; err != nil {
			return
		}
		result := Resource{
			Parent:   0,
			Fragment: "result",
			Name:     "结果页",
			Nav:      True(),
			Router:   False(),
			Icon:     "check-circle",
		}
		if err = tx.Create(&result).Error; err != nil {
			return
		}
		resultItems := []Resource{
			{
				Parent:   result.ID,
				Fragment: "success",
				Name:     "成功页",
				Nav:      True(),
				Router:   True(),
			},
			{
				Parent:   result.ID,
				Fragment: "fail",
				Name:     "失败页",
				Nav:      True(),
				Router:   True(),
			},
		}
		if err = tx.Create(&resultItems).Error; err != nil {
			return
		}
		exception := Resource{
			Parent:   0,
			Fragment: "exception",
			Name:     "异常页",
			Nav:      True(),
			Router:   False(),
			Icon:     "warning",
		}
		if err = tx.Create(&exception).Error; err != nil {
			return
		}
		exceptionItems := []Resource{
			{
				Parent:   exception.ID,
				Fragment: "403",
				Name:     "403",
				Nav:      True(),
				Router:   True(),
			},
			{
				Parent:   exception.ID,
				Fragment: "404",
				Name:     "404",
				Nav:      True(),
				Router:   True(),
			},
			{
				Parent:   exception.ID,
				Fragment: "500",
				Name:     "500",
				Nav:      True(),
				Router:   True(),
			},
		}
		if err = tx.Create(&exceptionItems).Error; err != nil {
			return
		}
		return
	}); err != nil {
		t.Error(err)
	}

	var resources []Resource
	if err := db.Where("strategy = ?", true).Find(&resources).Error; err != nil {
		t.Error(err)
	}
	r := make(map[string]uint64)
	for _, v := range resources {
		r[v.Fragment] = v.ID
	}
	data := []ResourceAclRel{
		{
			ResourceID: r["acl"],
			Path:       "acl",
			Mode:       "1",
		},
		{
			ResourceID: r["resource"],
			Path:       "resource",
			Mode:       "1",
		},
		{
			ResourceID: r["resource"],
			Path:       "acl",
			Mode:       "0",
		},
		{
			ResourceID: r["permission"],
			Path:       "permission",
			Mode:       "1",
		},
		{
			ResourceID: r["role"],
			Path:       "role",
			Mode:       "1",
		},
		{
			ResourceID: r["role"],
			Path:       "resource",
			Mode:       "0",
		},
		{
			ResourceID: r["role"],
			Path:       "permission",
			Mode:       "0",
		},
		{
			ResourceID: r["admin"],
			Path:       "admin",
			Mode:       "1",
		},
		{
			ResourceID: r["admin"],
			Path:       "role",
			Mode:       "0",
		},
		{
			ResourceID: r["admin"],
			Path:       "resource",
			Mode:       "0",
		},
		{
			ResourceID: r["admin"],
			Path:       "permission",
			Mode:       "0",
		},
	}
	if err := db.Create(&data).Error; err != nil {
		t.Error(err)
	}
}
