package model_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/weplanx/server/model"
	"testing"
)

func TestSetLogsetLogins(t *testing.T) {
	ctx := context.TODO()
	err := model.SetLogsetLogins(ctx, x.Db)
	assert.NoError(t, err)
}

func TestSetLogsetJobs(t *testing.T) {
	ctx := context.TODO()
	err := model.SetLogsetJobs(ctx, x.Db)
	assert.NoError(t, err)
}

func TestSetLogsetOperates(t *testing.T) {
	ctx := context.TODO()
	err := model.SetLogsetOperates(ctx, x.Db)
	assert.NoError(t, err)
}

func TestSetLogsetImessages(t *testing.T) {
	ctx := context.TODO()
	err := model.SetLogsetImessages(ctx, x.Db)
	assert.NoError(t, err)
}
