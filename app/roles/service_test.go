package roles

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"server/model"
	"testing"
)

func TestService_FindOneById(t *testing.T) {
	ctx := context.Background()
	var data model.Role
	if err := i.Db.Collection("roles").
		FindOne(ctx, bson.M{}).
		Decode(&data); err != nil {
		t.Error(err)
	}
	var role model.Role
	if err := service.FindOneById(ctx, data.ID, &role); err != nil {
		t.Error(err)
	}
	assert.Equal(t, data, role)
}

func TestService_FindNamesByIds(t *testing.T) {
	ctx := context.Background()
	var data model.Role
	if err := i.Db.Collection("roles").
		FindOne(ctx, bson.M{}).
		Decode(&data); err != nil {
		t.Error(err)
	}
	names, err := service.FindNamesByIds(ctx, []primitive.ObjectID{data.ID})
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, []string{data.Name}, names)
}

func TestService_FindByIds(t *testing.T) {
	ctx := context.Background()
	var data model.Role
	if err := i.Db.Collection("roles").
		FindOne(ctx, bson.M{}).
		Decode(&data); err != nil {
		t.Error(err)
	}
	var roles []model.Role
	if err := service.FindByIds(ctx, []primitive.ObjectID{data.ID}, &roles); err != nil {
		t.Error(err)
	}
	t.Log(roles)
}
