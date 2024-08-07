package index

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/weplanx/go/help"
	"github.com/weplanx/server/common"
	"net/http"
	"time"
)

func (x *Controller) GetRefreshCode(ctx context.Context, c *app.RequestContext) {
	claims := common.Claims(c)
	code, err := x.IndexX.GetRefreshCode(ctx, claims.UserId)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, M{
		"code": code,
	})
}

func (x *Service) GetRefreshCode(ctx context.Context, userId string) (code string, err error) {
	code = help.Uuid()
	x.Captcha.Create(ctx, userId, code, 15*time.Second)
	return
}
