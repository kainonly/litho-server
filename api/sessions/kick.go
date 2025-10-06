package sessions

import (
	"context"
	"server/common"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
)

type KickDto struct {
	Id string `json:"id" vd:"required"`
}

func (x *Controller) Kick(ctx context.Context, c *app.RequestContext) {
	var dto KickDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}
	user := common.GetIAM(c)
	n, err := x.SessionsX.Kick(ctx, user, dto.Id)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, utils.H{
		"count": n,
	})
}

func (x *Service) Kick(ctx context.Context, user *common.IAMUser, id string) (n int64, err error) {
	return x.RDb.Del(ctx, x.Key(id)).Result()
}
