package model_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/weplanx/server/model"
	"testing"
)

func TestSetProjects(t *testing.T) {
	ctx := context.TODO()
	err := model.SetProjects(ctx, x.Db)
	assert.NoError(t, err)
}
