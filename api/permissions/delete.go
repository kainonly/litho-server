package permissions

import (
	"context"
	"litho-api/common"
	"litho-api/model"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/kainonly/go/help"
)

func (x *Controller) Delete(ctx context.Context, c *app.RequestContext) {
	var dto common.DeleteDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	user := common.GetIAM(c)
	if err := x.PermissionsX.Delete(ctx, user, dto); err != nil {
		c.Error(err)
		return
	}

	common.SetTracking(c, Key, "delete", dto.IDs...)
	c.JSON(200, help.Ok())
}

func (x *Service) Delete(ctx context.Context, user *common.IAMUser, dto common.DeleteDto) (err error) {
	return x.Db.WithContext(ctx).Delete(model.Permission{}, dto.IDs).Error
}
