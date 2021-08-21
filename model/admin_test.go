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
			Roles:    []Role{{ID: 1}},
		},
		{
			Username:    "admin",
			Password:    password,
			Permissions: Array{"ACCESS_FINANCE_AUDIT"},
			Roles:       []Role{{ID: 2}},
		},
		{
			Username:  "test",
			Password:  password,
			Roles:     []Role{{ID: 3}, {ID: 4}},
			Resources: []Resource{{ID: 1}, {ID: 3}},
		},
	}
	if err := db.Create(&data).Error; err != nil {
		t.Error(err)
	}
}
