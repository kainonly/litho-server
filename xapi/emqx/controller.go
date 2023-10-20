package emqx

import (
	"context"
	"github.com/bytedance/gopkg/util/logger"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
)

type Controller struct {
	EmqxService *Service
}

type AuthDto struct {
	Identity string `json:"identity" vd:"required"`
	Token    string `json:"token" vd:"required"`
}

func (x *Controller) Auth(ctx context.Context, c *app.RequestContext) {
	var dto AuthDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	if err := x.EmqxService.Auth(ctx, dto); err != nil {
		c.Error(err)
		return
	}

	c.Status(204)
}

type AclDto struct {
	Identity string `json:"identity" vd:"mongodb"`
	Topic    string `json:"topic" vd:"required"`
}

func (x *Controller) Acl(ctx context.Context, c *app.RequestContext) {
	var dto AclDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	if err := x.EmqxService.Acl(ctx, dto); err != nil {
		logger.CtxErrorf(ctx, err.Error())
		c.JSON(200, utils.H{"result": "deny"})
		return
	}

	c.Status(204)
}

type BridgeDto struct {
	Client  string `json:"client" vd:"required"`
	Topic   string `json:"topic" vd:"required"`
	Payload M      `json:"payload" vd:"required,gt=0"`
}

func (x *Controller) Bridge(ctx context.Context, c *app.RequestContext) {
	var dto BridgeDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	if err := x.EmqxService.Bridge(ctx, dto); err != nil {
		c.Error(err)
		return
	}

	c.Status(204)
}
