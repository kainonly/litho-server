package roles

import (
	"context"
	"time"

	"server/common"
	"server/model"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/kainonly/go/help"
)

type UpdateDto struct {
	ID          string `json:"id" vd:"required"`
	OrgID       string `json:"org_id" vd:"required"`
	Name        string `json:"name" vd:"required"`
	Description string `json:"description" vd:"required"`
	Active      *bool  `json:"active" vd:"required"`
}

func (x *Controller) Update(ctx context.Context, c *app.RequestContext) {
	var dto UpdateDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	user := common.GetIAM(c)
	if err := x.RolesX.Update(ctx, user, dto); err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, help.Ok())
}

func (x *Service) Update(ctx context.Context, user *common.IAMUser, dto UpdateDto) (err error) {
	if err = x.Db.Model(model.Role{}).WithContext(ctx).
		Where(`id = ?`, dto.ID).
		Updates(utils.H{
			`updated_at`:  time.Now(),
			`org_id`:      dto.OrgID,
			`active`:      *dto.Active,
			`name`:        dto.Name,
			`description`: dto.Description,
		}).Error; err != nil {
		return
	}
	return
}
