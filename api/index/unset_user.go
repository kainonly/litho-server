package index

import (
	"context"
	"server/common"
	"server/model"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/kainonly/go/help"
)

type UnsetUserDto struct {
	Key string `json:"key" vd:"oneof='phone'"`
}

const IUnsetUser = "解绑资料"

func (x *Controller) UnsetUser(ctx context.Context, c *app.RequestContext) {
	var dto UnsetUserDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	user := common.GetIAM(c)
	if err := x.IndexX.UnsetUser(ctx, user, dto.Key); err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, help.Ok())
}

func (x *Service) UnsetUser(ctx context.Context, user *common.IAMUser, key string) (err error) {
	return x.Db.Model(model.User{}).WithContext(ctx).
		Where(`id = ?`, user.ID).
		Update(key, "").Error
}
