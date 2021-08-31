package model

import (
	"testing"
)

func TestRole(t *testing.T) {
	if err := db.AutoMigrate(&Role{}); err != nil {
		t.Error(err)
	}
	var resources []Resource
	if err := db.Find(&resources).Error; err != nil {
		t.Error(err)
	}
	role := []Role{
		{
			Name:        "超级管理员",
			Code:        "*",
			Description: "超级管理员拥有完整权限不能编辑，若不使用可以禁用该权限",
		},
		{
			Name:        "管理员",
			Code:        "admin",
			Description: "默认",
			Resources:   resources,
		},
		{
			Name:        "分析员",
			Code:        "analysis",
			Description: "默认",
			Resources: []Resource{
				{Path: "dashboard"},
				{Path: "dashboard/analysis"},
				{Path: "dashboard/monitor"},
				{Path: "dashboard/workbench"},
			},
		},
		{
			Name:        "成员",
			Code:        "staff",
			Description: "默认",
			Resources: []Resource{
				{Path: "exception"},
				{Path: "exception/403"},
				{Path: "exception/404"},
				{Path: "exception/500"},
			},
		},
	}
	if err := db.Create(role).Error; err != nil {
		t.Error(err)
	}
}
