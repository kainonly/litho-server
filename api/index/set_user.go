package index

import (
	"context"
	"server/common"
	"server/model"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/kainonly/go/help"
)

type SetUserDto struct {
	Key    string `json:"key" vd:"oneof='email' 'name' 'avatar'"`
	Email  string `json:"email" vd:"required_if=Key 'Email',omitempty,email"`
	Name   string `json:"name" vd:"required_if=Key 'Name'"`
	Avatar string `json:"avatar" vd:"required_if=Key 'Avatar'"`
}

func (x *Controller) SetUser(ctx context.Context, c *app.RequestContext) {
	var dto SetUserDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	user := common.GetIAM(c)
	if err := x.IndexX.SetUser(ctx, user, dto); err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, help.Ok())
}

func (x *Service) SetUser(ctx context.Context, user *common.IAMUser, dto SetUserDto) (err error) {
	data := utils.H{}
	switch dto.Key {
	case "phone":
		data["phone"] = dto.Email
		break
	case "name":
		data["name"] = dto.Name
		break
	case "avatar":
		data["avatar"] = dto.Avatar
		break
	}

	if err = x.Db.Model(model.User{}).WithContext(ctx).
		Where(`id = ?`, user.ID).
		Updates(data).Error; err != nil {
		return
	}

	return
}
