package lark

import (
	"context"
	"github.com/bytedance/sonic"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/weplanx/go/help"
)

type ChallengeDto struct {
	Encrypt string `json:"encrypt" vd:"required"`
}

func (x *Controller) Challenge(ctx context.Context, c *app.RequestContext) {
	var dto ChallengeDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}
	raw, err := x.LarkX.Decrypt(dto.Encrypt, x.V.LarkEncryptKey)
	if err != nil {
		c.Error(err)
		return
	}
	var data struct {
		Challenge string `json:"challenge"`
		Token     string `json:"token"`
		Type      string `json:"type"`
	}
	if err = sonic.UnmarshalString(raw, &data); err != nil {
		c.Error(err)
		return
	}
	if data.Token != x.V.LarkVerificationToken {
		c.Error(help.E(
			"lark.VerificationTokenNotMatch",
			"the local configuration token does not match the authentication token"),
		)
		return
	}

	c.JSON(200, utils.H{
		"challenge": data.Challenge,
	})
}
