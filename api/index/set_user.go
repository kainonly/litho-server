package index

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/weplanx/server/common"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"reflect"
	"strings"
)

type SetUserDto struct {
	Key    string `json:"key" vd:"oneof='email' 'name' 'avatar'"`
	Email  string `json:"email" vd:"required_if=Key 'Email',omitempty,email"`
	Name   string `json:"name" vd:"required_if=Key 'Name'"`
	Avatar string `json:"avatar" vd:"required_if=Key 'Avatar'"`
}

func (x *Controller) SetUser(ctx context.Context, c *app.RequestContext) {
	var dto SetUserDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	data := make(M)
	key := strings.ToUpper(dto.Key[:1]) + dto.Key[1:]
	data[dto.Key] = reflect.ValueOf(dto).
		FieldByName(key).
		Interface()

	claims := common.Claims(c)
	if _, err := x.IndexX.SetUser(ctx, claims.UserId, bson.M{"$set": data}); err != nil {
		c.Error(err)
		return
	}

	c.Status(204)
}

func (x *Service) SetUser(ctx context.Context, userId string, update bson.M) (result interface{}, err error) {
	id, _ := primitive.ObjectIDFromHex(userId)

	if result, err = x.Db.Collection("users").
		UpdateByID(ctx, id, update); err != nil {
		return
	}

	key := fmt.Sprintf(`users:%s`, userId)
	if _, err = x.RDb.Del(ctx, key).Result(); err != nil {
		return
	}

	return
}
