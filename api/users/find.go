package users

import (
	"context"

	"server/common"
	"server/model"

	"github.com/cloudwego/hertz/pkg/app"
)

type FindDto struct {
	common.FindDto
	Active *bool  `query:"active,omitempty"`
	Email  string `query:"email,omitempty"`
	Phone  string `query:"phone,omitempty"`
}

func (x *Controller) Find(ctx context.Context, c *app.RequestContext) {
	var dto FindDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.JSON(400, map[string]any{"error": err.Error()})
		return
	}
	user := c.MustGet("user").(*common.IAMUser)
	data, err := x.UsersX.Find(ctx, user, dto)
	if err != nil {
		c.JSON(500, map[string]any{"error": err.Error()})
		return
	}
	c.JSON(200, data)
}

func (x *Service) Find(ctx context.Context, user *common.IAMUser, dto FindDto) (data []model.User, err error) {
	ctx = common.SetPipe(ctx, common.NewFindPipe())
	do := x.Db.Model(&model.User{}).Omit("password")
	if dto.Active != nil {
		do = do.Where("active = ?", *dto.Active)
	}
	if dto.Email != "" {
		do = do.Where("email = ?", dto.Email)
	}
	if dto.Phone != "" {
		do = do.Where("phone = ?", dto.Phone)
	}
	if dto.Q != "" {
		do = do.Where("name LIKE ? OR email LIKE ? OR phone LIKE ?",
			dto.GetKeyword(), dto.GetKeyword(), dto.GetKeyword())
	}
	err = dto.FindDto.Find(ctx, do, &data)
	return
}
