package teams

import (
	"context"
	"server/common"
	"server/model"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/kainonly/go/help"
)

type CreateDto struct {
	ID   string `json:"-"`
	Key  string `json:"key" vd:"required"`
	Name string `json:"name" vd:"required"`
}

func (x *Controller) Create(ctx context.Context, c *app.RequestContext) {
	var dto CreateDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	dto.ID = help.SID()
	user := common.GetIAM(c)
	if err := x.TeamsX.Create(ctx, user, dto); err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, help.Ok())
}

func (x *Service) Create(ctx context.Context, user *common.IAMUser, dto CreateDto) (err error) {
	data := &model.Team{
		ID:   dto.ID,
		Key:  dto.Key,
		Name: dto.Name,
	}
	return x.Db.WithContext(ctx).Create(data).Error
}
