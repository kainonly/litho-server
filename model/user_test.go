package model_test

import (
	"github.com/alexedwards/argon2id"
	"github.com/stretchr/testify/assert"
	"github.com/weplanx/server/model"
	"testing"
)

func TestCreateUser(t *testing.T) {
	err := db.Migrator().DropTable(model.User{})
	assert.NoError(t, err)
	err = db.AutoMigrate(model.User{})
	assert.NoError(t, err)

	hash, err := argon2id.CreateHash("pass@VAN1234", argon2id.DefaultParams)
	assert.NoError(t, err)

	err = db.Create(&model.User{
		Email:    "zhangtqx@vip.qq.com",
		Password: hash,
	}).Error
	assert.NoError(t, err)

}
