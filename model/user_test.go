package model_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/weplanx/server/model"
	"testing"
)

func TestSetupUser(t *testing.T) {
	ctx := context.TODO()
	err := model.SetupUser(ctx, x.Db)
	assert.NoError(t, err)
}

func TestCreateUser(t *testing.T) {
	_, err := x.Db.Collection("users").InsertOne(
		context.TODO(),
		model.NewUser("zhangtqx@qq.com", "pass@VAN1234"),
	)
	assert.NoError(t, err)
}
