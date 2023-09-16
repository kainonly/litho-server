package model_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/weplanx/server/model"
	"testing"
)

func TestSetQueues(t *testing.T) {
	ctx := context.TODO()
	err := model.SetQueues(ctx, x.Db)
	assert.NoError(t, err)
}
