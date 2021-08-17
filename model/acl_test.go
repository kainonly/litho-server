package model

import "testing"

func TestAcl(t *testing.T) {
	if err := db.AutoMigrate(&Acl{}); err != nil {
		t.Error(err)
	}
}
