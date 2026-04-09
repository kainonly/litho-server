package departments

import (
	"context"
	"server/model"

	"server/common"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/kainonly/go/help"
)

const IDelete = "删除"

func (x *Controller) Delete(ctx context.Context, c *app.RequestContext) {
	var dto common.DeleteDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	user := common.GetIAM(c)
	if err := x.DepartmentsX.Delete(ctx, user, dto); err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, help.Ok())
}

func (x *Service) Delete(ctx context.Context, user *common.IAMUser, dto common.DeleteDto) (err error) {
	return x.Db.WithContext(ctx).Delete(model.Department{}, dto.IDs).Error
}
