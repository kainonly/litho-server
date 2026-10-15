package orgs

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
	if err := x.OrgsX.Delete(ctx, user, dto); err != nil {
		c.Error(err)
		return
	}

	common.SetTracking(c, Key, "delete", dto.IDs...)
	c.JSON(200, help.Ok())
}

func (x *Service) Delete(ctx context.Context, user *common.IAMUser, dto common.DeleteDto) (err error) {
	if err = x.CheckAccess(ctx, user); err != nil {
		return
	}

	var exists int64
	if err = x.Db.Model(model.UserOrg{}).WithContext(ctx).
		Where(`org_id in (?)`, dto.IDs).
		Count(&exists).Error; err != nil {
		return
	}
	if exists != 0 {
		return help.E(0, "禁止删除，该组织已存在用户")
	}
	return x.Db.WithContext(ctx).Delete(model.Org{}, dto.IDs).Error
}
