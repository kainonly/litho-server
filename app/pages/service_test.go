package pages

import (
	"api/model"
	"context"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
)

func TestService_FindOneById(t *testing.T) {
	ctx := context.Background()
	var data model.Page
	if err := i.Db.Collection("pages").
		FindOne(ctx, bson.M{}).
		Decode(&data); err != nil {
		t.Error(err)
	}
	var page model.Page
	if err := service.FindOneById(ctx, data.ID, &page); err != nil {
		t.Error(err)
	}
	assert.Equal(t, data, page)
}

func TestService_Navs(t *testing.T) {
	ctx := context.Background()
	cursor, err := i.Db.Collection("roles").
		Find(ctx, bson.M{})
	if err != nil {
		t.Error(err)
	}
	var roles []model.Role
	if err = cursor.All(ctx, &roles); err != nil {
		t.Error(err)
	}
	assert.NotEmpty(t, roles)
	navs, err := service.Navs(ctx, []model.Role{})
	if err != nil {
		t.Error(err)
	}
	t.Log(navs)
}
