package model_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/weplanx/server/model"
	"testing"
)

func TestCreateRole(t *testing.T) {
	err := db.Migrator().DropTable(model.Role{})
	assert.NoError(t, err)
	err = db.AutoMigrate(model.Role{})
	assert.NoError(t, err)

	err = db.Create(&model.Role{
		Name:        "超级管理员",
		Description: "系统",
		//Pages: map[string]int64{
		//	"admin": 1,
		//},
	}).Error
	assert.NoError(t, err)
}
