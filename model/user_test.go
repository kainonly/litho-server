package model_test

import (
	"context"
	"github.com/alexedwards/argon2id"
	"github.com/stretchr/testify/assert"
	"github.com/weplanx/server/model"
	"testing"
)

func TestCreateUser(t *testing.T) {
	hash, err := argon2id.CreateHash("pass@VAN1234", argon2id.DefaultParams)
	assert.NoError(t, err)
	_, err = db.Collection("users").InsertOne(
		context.TODO(),
		model.NewUser("zhangtqx@vip.qq.com", hash),
	)
	assert.NoError(t, err)
}
