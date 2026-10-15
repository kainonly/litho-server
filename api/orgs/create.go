package orgs

import (
	"context"
	"litho-api/common"
	"litho-api/model"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/kainonly/go/help"
)

type CreateDto struct {
	ID      string `json:"-"`
	PID     string `json:"pid" vd:"required"`
	Code    string `json:"code"`
	Name    string `json:"name" vd:"required"`
	Type    *int16 `json:"type" vd:"required"`
	Company *bool  `json:"company" vd:"required"`
	Status  *bool  `json:"status" vd:"required"`
}

func (x *Controller) Create(ctx context.Context, c *app.RequestContext) {
	var dto CreateDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	user := common.GetIAM(c)
	dto.ID = help.SID()
	if err := x.OrgsX.Create(ctx, user, dto); err != nil {
		c.Error(err)
		return
	}

	common.SetTracking(c, Key, "create", dto.ID)
	c.JSON(200, help.Ok())
}

func (x *Service) Create(ctx context.Context, user *common.IAMUser, dto CreateDto) (err error) {
	if err = x.CheckAccess(ctx, user); err != nil {
		return
	}

	data := &model.Org{
		ID:        dto.ID,
		Status:    dto.Status,
		Pid:       dto.PID,
		Code:      dto.Code,
		Name:      dto.Name,
		Type:      dto.Type,
		Company:   *dto.Company,
		ContactID: `0`,
	}
	return x.Db.WithContext(ctx).Create(data).Error
}
