package model

import (
	"testing"
)

func TestRole(t *testing.T) {
	if err := db.AutoMigrate(&Role{}); err != nil {
		t.Error(err)
	}
	var resources []Resource
	if err := db.Select("id").Find(&resources).Error; err != nil {
		t.Error(err)
	}
	if err := db.Create(&Role{
		Key:      "*",
		Name:     "超级管理员",
		Resource: resources,
	}); err != nil {
		t.Error(err)
	}
}
