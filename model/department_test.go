package model_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/weplanx/server/model"
	"testing"
)

func TestCreateDepartment(t *testing.T) {
	err := db.Migrator().DropTable(model.Department{})
	assert.NoError(t, err)
	err = db.AutoMigrate(model.Department{})
	assert.NoError(t, err)

	err = db.Create(&model.Department{Name: "运营"}).Error
	assert.NoError(t, err)
}
