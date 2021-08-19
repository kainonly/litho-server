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
	if err := db.Create(&Admin{
		Username: "kain",
		Password: password,
	}).Error; err != nil {
		t.Error(err)
	}
}
