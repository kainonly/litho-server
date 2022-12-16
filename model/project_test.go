package model_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/weplanx/server/model"
	"testing"
	"time"
)

func TestCreateProject(t *testing.T) {
	var err error
	_, err = db.Collection("projects").InsertOne(
		context.TODO(),
		model.Project{
			Name:        "默认",
			Namespace:   "default",
			Entry:       []string{},
			Expire:      0,
			Status:      true,
			CreatedTime: time.Now(),
			UpdatedTime: time.Now(),
		},
	)
	assert.NoError(t, err)
}
