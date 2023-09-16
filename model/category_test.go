package model_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/weplanx/server/model"
	"testing"
)

func TestSetCategories(t *testing.T) {
	ctx := context.TODO()
	err := model.SetCategories(ctx, x.Db)
	assert.NoError(t, err)
}
