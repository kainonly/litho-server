package model

import (
	"testing"
)

func TestRole(t *testing.T) {
	if err := db.AutoMigrate(&Role{}); err != nil {
		t.Error(err)
	}
	if err := db.Create(&Role{
		Name:        "超级管理员",
		Description: "超级管理员权限不能删除，若不使用可以禁用该权限",
	}).Error; err != nil {
		t.Error(err)
	}
}
