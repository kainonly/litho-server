package model_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/weplanx/server/model"
	"testing"
)

func TestCreateProject(t *testing.T) {
	err := db.Migrator().DropTable(model.Project{})
	assert.NoError(t, err)
	err = db.AutoMigrate(model.Project{})
	assert.NoError(t, err)

	err = db.Create(&model.Project{
		Name:      "默认",
		Namespace: "default",
	}).Error
	assert.NoError(t, err)
}
