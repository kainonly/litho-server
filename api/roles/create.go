package roles

import (
	"context"

	"server/common"
	"server/model"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/kainonly/go/help"
)

type CreateDto struct {
	ID          string `json:"-"`
	OrgID       string `json:"org_id" vd:"required"`
	Name        string `json:"name" vd:"required"`
	Description string `json:"description" vd:"required"`
	Sort        int16  `json:"sort"`
	Active      *bool  `json:"active" vd:"required"`
}

func (x *Controller) Create(ctx context.Context, c *app.RequestContext) {
	var dto CreateDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	user := common.GetIAM(c)
	dto.ID = help.SID()
	if err := x.RolesX.Create(ctx, user, dto); err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, help.Ok())
}

func (x *Service) Create(ctx context.Context, user *common.IAMUser, dto CreateDto) (err error) {
	data := model.Role{
		ID:          dto.ID,
		OrgID:       dto.OrgID,
		Sort:        dto.Sort,
		Active:      *dto.Active,
		Name:        dto.Name,
		Description: dto.Description,
	}
	if err = x.Db.WithContext(ctx).
		Create(&data).Error; err != nil {
		return
	}
	return
}
