package model_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/weplanx/server/model"
	"testing"
)

func TestSetBuilders(t *testing.T) {
	ctx := context.TODO()
	err := model.SetBuilders(ctx, x.Db)
	assert.NoError(t, err)
}
