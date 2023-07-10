package index

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"net/http"
	"time"
)

type Controller struct {
	IndexService *Service
}

func (x *Controller) Ping(ctx context.Context, c *app.RequestContext) {
	c.JSON(http.StatusOK, utils.H{
		"ip":   c.ClientIP(),
		"time": time.Now(),
	})
}

type LoginDto struct {
	Email    string `json:"email,required" vd:"email($)"`
	Password string `json:"password,required" vd:"len($)>=8"`
}

//func (x *Controller) Login(ctx context.Context, c *app.RequestContext) {
//	var dto LoginDto
//	if err := c.BindAndValidate(&dto); err != nil {
//		c.Error(err)
//		return
//	}
//
//	ts, err := x.Service.Login(ctx, dto.Email, dto.Password, logdata)
//	if err != nil {
//		c.Error(err)
//		return
//	}
//
//	c.SetCookie("access_token", ts, 0, "/", "", protocol.CookieSameSiteLaxMode, true, true)
//	c.JSON(200, utils.H{
//		"code":    0,
//		"message": "ok",
//	})
//}
