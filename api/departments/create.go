package departments

import (
	"context"
	"server/common"
	"server/model"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/kainonly/go/help"
)

type CreateDto struct {
	ID     string `json:"-"`
	Type   *int16 `json:"type" vd:"required"`
	Name   string `json:"name" vd:"required"`
	Status *bool  `json:"status" vd:"required"`
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
	if err := x.DepartmentsX.Create(ctx, user, dto); err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, help.Ok())
}

func (x *Service) Create(ctx context.Context, user *common.IAMUser, dto CreateDto) (err error) {
	data := model.Department{
		ID:     dto.ID,
		Status: dto.Status,
		Type:   dto.Type,
		Name:   dto.Name,
	}
	if err = x.Db.WithContext(ctx).
		Create(&data).Error; err != nil {
		return
	}
	return
}
