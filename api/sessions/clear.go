package sessions

import (
	"context"
	"server/common"

	"github.com/cloudwego/hertz/pkg/app"
)

func (x *Controller) Clear(ctx context.Context, c *app.RequestContext) {
	user := common.GetIAM(c)
	n, err := x.SessionsX.Clear(ctx, user)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, common.M{
		"count": n,
	})
}

func (x *Service) Clear(ctx context.Context, user *common.IAMUser) (n int64, err error) {
	var keys []string
	x.Scan(ctx, func(key string) {
		keys = append(keys, key)
	})
	if len(keys) != 0 {
		return x.RDb.Del(ctx, keys...).Result()
	}
	return
}
