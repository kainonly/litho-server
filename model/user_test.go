package model_test

import (
	"github.com/alexedwards/argon2id"
	"github.com/go-faker/faker/v4"
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
		Username: "weplanx",
		Password: hash,
		Email:    "zhangtqx@vip.qq.com",
		Roles:    []int64{1},
	}).Error
	assert.NoError(t, err)

	var mock []model.User
	for i := 0; i < 10; i++ {
		p, _ := faker.RandomInt(1, 10, 5)
		roles := []int64{}
		for _, v := range p {
			roles = append(roles, int64(v))
		}
		mock = append(mock, model.User{
			Username: faker.Username(),
			Password: hash,
			Email:    faker.Email(),
			Roles:    roles,
		})
	}

	err = db.Create(mock).Error
	assert.NoError(t, err)
}
