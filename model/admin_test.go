package model

import (
	"github.com/kainonly/go-bit/hash"
	"testing"
)

func TestAdmin(t *testing.T) {
	if err := db.AutoMigrate(&Admin{}); err != nil {
		t.Error(err)
	}
	password, _ := hash.Make("pass@VAN1234")
	data := []Admin{
		{
			Username: "kain",
			Password: password,
			Roles: []Role{
				{Code: "*"},
			},
		},
		{
			Username: "admin",
			Password: password,
			Roles: []Role{
				{Code: "admin"},
			},
		},
		{
			Username: "test",
			Password: password,
			Roles: []Role{
				{Code: "analysis"},
				{Code: "staff"},
			},
			Resources: []Resource{
				{Path: "form"},
				{Path: "form/basic"},
				{Path: "form/step"},
				{Path: "form/advanced"},
			},
		},
	}
	if err := db.Create(&data).Error; err != nil {
		t.Error(err)
	}
}
