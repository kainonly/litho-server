package emqx

import (
	"context"
	"github.com/bytedance/gopkg/util/logger"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Controller struct {
	EmqxService *Service
}

type AuthDto struct {
	Identity primitive.ObjectID `json:"identity,required"`
	Token    string             `json:"token,required"`
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
	Identity primitive.ObjectID `json:"identity,required"`
	Topic    string             `json:"topic,required"`
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
