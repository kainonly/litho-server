package roles

import (
	"context"

	"server/common"
	"server/model"

	"github.com/cloudwego/hertz/pkg/app"
)

type UpdateDto struct {
	ID          string  `json:"id" vd:"required"`
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	Sort        *int16  `json:"sort,omitempty"`
	Active      *bool   `json:"active,omitempty"`
}

func (x *Controller) Update(ctx context.Context, c *app.RequestContext) {
	var dto UpdateDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.JSON(400, map[string]any{"error": err.Error()})
		return
	}
	user := c.MustGet("user").(*common.IAMUser)
	if err := x.RolesX.Update(ctx, user, dto); err != nil {
		c.JSON(500, map[string]any{"error": err.Error()})
		return
	}
	c.JSON(200, map[string]any{"ok": true})
}

func (x *Service) Update(ctx context.Context, user *common.IAMUser, dto UpdateDto) (err error) {
	data := make(map[string]any)
	if dto.Name != nil {
		data["name"] = *dto.Name
	}
	if dto.Description != nil {
		data["description"] = *dto.Description
	}
	if dto.Sort != nil {
		data["sort"] = *dto.Sort
	}
	if dto.Active != nil {
		data["active"] = *dto.Active
	}
	return x.Db.Model(&model.Role{}).Where(`id = ?`, dto.ID).Updates(data).Error
}
