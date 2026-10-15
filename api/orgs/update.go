package orgs

import (
	"context"
	"litho-api/common"
	"litho-api/model"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/kainonly/go/help"
)

type UpdateDto struct {
	ID      string `json:"id"`
	PID     string `json:"pid" vd:"required"`
	Code    string `json:"code"`
	Name    string `json:"name" vd:"required"`
	Type    *int16 `json:"type" vd:"required"`
	Company *bool  `json:"company" vd:"required"`
	Status  *bool  `json:"status" vd:"required"`
}

func (x *Controller) Update(ctx context.Context, c *app.RequestContext) {
	var dto UpdateDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	user := common.GetIAM(c)
	if err := x.OrgsX.Update(ctx, user, dto); err != nil {
		c.Error(err)
		return
	}

	common.SetTracking(c, Key, "update", dto.ID)
	c.JSON(200, help.Ok())
}

func (x *Service) Update(ctx context.Context, user *common.IAMUser, dto UpdateDto) (err error) {
	if err = x.CheckAccess(ctx, user); err != nil {
		return
	}

	updates := M{
		"update_time": time.Now(),
		"pid":         dto.PID,
		"code":        dto.Code,
		"name":        dto.Name,
		"type":        *dto.Type,
		"company":     *dto.Company,
		"status":      *dto.Status,
	}
	return x.Db.Model(model.Org{}).WithContext(ctx).
		Where(`id = ?`, dto.ID).
		Updates(updates).Error
}
