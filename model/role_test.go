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
			Description: "超级管理员拥有完整权限不能编辑，若不使用可以禁用该权限",
		},
		{
			Name:        "管理员",
			Description: "默认",
			Resources:   resources,
			Permissions: []interface{}{"ACCESS_FINANCE_AUDIT"},
		},
		{
			Name:        "权限管理员",
			Description: "默认",
			Resources:   []Resource{{ID: 1}, {ID: 5}},
		},
		{
			Name:        "成员管理员",
			Description: "默认",
			Resources:   []Resource{{ID: 1}, {ID: 6}},
		},
	}
	if err := db.Create(role).Error; err != nil {
		t.Error(err)
	}
}
