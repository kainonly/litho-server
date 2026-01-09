package roles

import (
	"context"

	"server/common"
	"server/model"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/kainonly/go/snowflake"
)

type CreateDto struct {
	Name        string `json:"name" vd:"required"`
	Description string `json:"description"`
	Sort        int16  `json:"sort"`
}

func (x *Controller) Create(ctx context.Context, c *app.RequestContext) {
	var dto CreateDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.JSON(400, map[string]any{"error": err.Error()})
		return
	}
	user := c.MustGet("user").(*common.IAMUser)
	if err := x.RolesX.Create(ctx, user, dto); err != nil {
		c.JSON(500, map[string]any{"error": err.Error()})
		return
	}
	c.JSON(201, map[string]any{"ok": true})
}

func (x *Service) Create(ctx context.Context, user *common.IAMUser, dto CreateDto) (err error) {
	data := model.Role{
		ID:          snowflake.GenerateString(),
		OrgID:       user.OrgID,
		Name:        dto.Name,
		Description: dto.Description,
		Sort:        dto.Sort,
		Active:      true,
	}
	return x.Db.Create(&data).Error
}
