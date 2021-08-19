package model

import "testing"

func TestPermission(t *testing.T) {
	if err := db.AutoMigrate(&Permission{}); err != nil {
		t.Error(err)
	}
	data := []Permission{
		{
			Code: "ACCESS_FINANCE_AUDIT",
			Name: "允许财务审计",
		},
	}
	if err := db.Create(&data).Error; err != nil {
		t.Error(err)
	}
}
