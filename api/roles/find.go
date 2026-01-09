package roles

import (
	"context"

	"server/common"
	"server/model"

	"github.com/cloudwego/hertz/pkg/app"
)

func (x *Controller) Find(ctx context.Context, c *app.RequestContext) {
	var dto common.FindDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.JSON(400, map[string]any{"error": err.Error()})
		return
	}
	user := c.MustGet("user").(*common.IAMUser)
	data, err := x.RolesX.Find(ctx, user, dto)
	if err != nil {
		c.JSON(500, map[string]any{"error": err.Error()})
		return
	}
	c.JSON(200, data)
}

func (x *Service) Find(ctx context.Context, user *common.IAMUser, dto common.FindDto) (data []model.Role, err error) {
	ctx = common.SetPipe(ctx, common.NewFindPipe())
	do := x.Db.Model(&model.Role{}).Where(`org_id = ?`, user.OrgID)
	if dto.Q != "" {
		do = do.Where(`name ILIKE ?`, dto.GetKeyword())
	}
	err = dto.Find(ctx, do, &data)
	return
}
