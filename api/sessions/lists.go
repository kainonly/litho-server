package sessions

import (
	"context"
	"server/common"
	"server/model"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
)

func (x *Controller) Lists(ctx context.Context, c *app.RequestContext) {
	user := common.GetIAM(c)
	result, err := x.SessionsX.Lists(ctx, user)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, result)
}

type ListsResult struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
	Sessions int32  `json:"sessions"`
}

func (x *Service) Lists(ctx context.Context, user *common.IAMUser) (result []ListsResult, err error) {
	ids := make([]string, 0)
	x.Scan(ctx, func(key string) {
		v := strings.Replace(key, x.Key(""), "", -1)
		ids = append(ids, v)
	})
	result = make([]ListsResult, len(ids))
	if err = x.Db.Model(model.User{}).WithContext(ctx).
		Select([]string{"id", "email", "sessions"}).
		Where(`id in (?)`, ids).
		Find(&result).Error; err != nil {
		return
	}
	return
}
