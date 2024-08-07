package index

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/weplanx/server/common"
	"go.mongodb.org/mongo-driver/bson"
)

type UnsetUserDto struct {
	Key string `path:"key" vd:"oneof='phone' 'totp' 'lark'"`
}

func (x *Controller) UnsetUser(ctx context.Context, c *app.RequestContext) {
	var dto UnsetUserDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	claims := common.Claims(c)
	_, err := x.IndexX.SetUser(ctx, claims.UserId, bson.M{
		"$unset": bson.M{dto.Key: 1},
	})
	if err != nil {
		c.Error(err)
		return
	}

	c.Status(204)
}
