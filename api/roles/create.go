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
	Name        string `json:"name" vd:"required"`
	Description string `json:"description" vd:"required"`
	Sort        int16  `json:"sort"`
	Status      *bool  `json:"status" vd:"required"`
}

const ICreate = "新增"

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
		Sort:        dto.Sort,
		Status:      dto.Status,
		Name:        dto.Name,
		Description: dto.Description,
		Strategy: common.RoleStrategy{
			Navs:        make([]string, 0),
			Routes:      make([]string, 0),
			Permissions: make([]string, 0),
		},
	}
	if err = x.Db.WithContext(ctx).
		Create(&data).Error; err != nil {
		return
	}
	return
}
