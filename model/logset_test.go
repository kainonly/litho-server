package model_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/weplanx/server/model"
	"testing"
)

func TestSetupLogsetLogined(t *testing.T) {
	ctx := context.TODO()
	err := model.SetupLogsetLogined(ctx, x.Db)
	assert.NoError(t, err)
}

func TestSetupLogsetJobs(t *testing.T) {
	ctx := context.TODO()
	err := model.SetupLogsetJobs(ctx, x.Db)
	assert.NoError(t, err)
}

func TestSetupLogsetOperates(t *testing.T) {
	ctx := context.TODO()
	err := model.SetupLogsetOperates(ctx, x.Db)
	assert.NoError(t, err)
}
