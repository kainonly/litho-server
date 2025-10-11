package index

import (
	"context"
	"server/common"
	"server/model"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/errors"
	"github.com/kainonly/go/help"
	"github.com/kainonly/go/passlib"
)

type SetUserPasswordDto struct {
	Old      string `json:"old"`
	Password string `json:"password" vd:"min=8"`
}

func (x *Controller) SetUserPassword(ctx context.Context, c *app.RequestContext) {
	var dto SetUserPasswordDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	user := common.GetIAM(c)
	if err := x.IndexX.SetUserPassword(ctx, user, dto); err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, help.Ok())
}

func (x *Service) SetUserPassword(ctx context.Context, user *common.IAMUser, dto SetUserPasswordDto) (err error) {
	var data model.User
	if err = x.Db.Model(model.User{}).WithContext(ctx).
		Where(`id = ?`, user.ID).
		Take(&data).Error; err != nil {
		return
	}

	if dto.Old == "" && data.Password == "" {
		var hash string
		if hash, err = passlib.Hash(dto.Password); err != nil {
			return
		}
		return x.Db.Model(model.User{}).WithContext(ctx).
			Where(`id = ?`, user.ID).
			Update("password", hash).Error
	}

	if err = passlib.Verify(dto.Old, data.Password); err != nil {
		if err == passlib.ErrNotMatch {
			err = errors.NewPublic(passlib.ErrNotMatch.Error())
			return
		}
	}

	var hash string
	if hash, err = passlib.Hash(dto.Password); err != nil {
		return
	}

	return x.Db.Model(model.User{}).WithContext(ctx).
		Where(`id = ?`, user.ID).
		Update("password", hash).Error
}
