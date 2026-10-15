package sessions

import (
	"context"
	"litho-api/common"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
)

func (x *Controller) Clear(ctx context.Context, c *app.RequestContext) {
	user := common.GetIAM(c)
	n, err := x.SessionsX.Clear(ctx, user)
	if err != nil {
		c.Error(err)
		return
	}

	common.SetTracking(c, Key, "clear", "0")
	c.JSON(200, utils.H{
		"count": n,
	})
}

func (x *Service) Clear(ctx context.Context, user *common.IAMUser) (n int64, err error) {
	if err = x.CheckAccess(ctx, user); err != nil {
		return
	}

	var keys []string
	x.Scan(ctx, func(key string) {
		keys = append(keys, key)
	})
	if len(keys) != 0 {
		return x.RDb.Del(ctx, keys...).Result()
	}
	return
}
