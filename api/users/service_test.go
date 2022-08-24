package users_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/weplanx/api/model"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
)

func TestService_GetActived(t *testing.T) {
	var err error
	var user model.User
	err = x.Db.Collection("users").
		FindOne(context.TODO(), bson.M{}).
		Decode(&user)
	assert.NoError(t, err)

	var data model.User
	data, err = x.UsersService.GetActived(context.TODO(), user.ID.Hex())
	assert.NoError(t, err)
	t.Log(data)
}
