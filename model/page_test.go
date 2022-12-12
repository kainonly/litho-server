package model_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/weplanx/server/model"
	"testing"
)

func TestCreatePage(t *testing.T) {
	err := db.Migrator().DropTable(model.Page{})
	assert.NoError(t, err)
	err = db.AutoMigrate(model.Page{})
	assert.NoError(t, err)

	err = db.Create(&model.Page{
		Name: "测试",
	}).Error
	assert.NoError(t, err)
}
