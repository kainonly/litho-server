package emqx

import (
	"context"
	"fmt"
	"github.com/bytedance/gopkg/util/logger"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/errors"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/weplanx/server/model"
	"go.mongodb.org/mongo-driver/bson"
	"strings"
)

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

func (x *Service) Acl(ctx context.Context, dto AclDto) (err error) {
	deny := true
	topic := strings.Split(dto.Topic, "/")
	msg := fmt.Sprintf(`The user [%s] is not authorized for this topic [%s]`,
		dto.Identity, dto.Topic)
	if !(len(topic) >= 2 && topic[1] == dto.Identity) {
		return errors.NewPublic(msg)
	}
	var data model.Imessage
	if err = x.Db.Collection("imessages").
		FindOne(ctx, bson.M{"topic": topic[0]}).
		Decode(&data); err != nil {
		return
	}
	for _, pid := range data.Projects {
		if pid.Hex() == dto.Identity {
			deny = false
			break
		}
	}
	if deny {
		return errors.NewPublic(msg)
	}
	return
}
