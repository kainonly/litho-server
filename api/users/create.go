package users

import (
	"context"
	"litho-api/common"
	"litho-api/model"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/kainonly/go/help"
	"github.com/kainonly/go/passlib"
	"gorm.io/gorm"
)

type CreateDto struct {
	ID       string `json:"-"`
	OrgID    string `json:"org_id" vd:"required"`
	Phone    string `json:"phone" vd:"required"`
	Name     string `json:"name" vd:"required"`
	Password string `json:"password" vd:"required"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
	Status   *bool  `json:"status"`
}

func (x *Controller) Create(ctx context.Context, c *app.RequestContext) {
	var dto CreateDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	dto.ID = help.SID()
	user := common.GetIAM(c)
	if err := x.UsersX.Create(ctx, user, dto); err != nil {
		c.Error(err)
		return
	}

	common.SetTracking(c, Key, "create", dto.ID)
	c.JSON(200, help.Ok())
}

func (x *Service) Create(ctx context.Context, user *common.IAMUser, dto CreateDto) (err error) {
	if err = x.CheckAccess(ctx, user); err != nil {
		return
	}

	status := true
	if dto.Status != nil {
		status = *dto.Status
	}
	data := &model.User{
		ID:     dto.ID,
		Status: &status,
		Name:   dto.Name,
		Phone:  dto.Phone,
		Email:  dto.Email,
		Avatar: dto.Avatar,
	}
	data.Password, _ = passlib.Hash(dto.Password)
	if err = x.Db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err = tx.Create(data).Error; err != nil {
			return err
		}
		return tx.Create(&model.UserOrg{
			OrgID:  dto.OrgID,
			UserID: dto.ID,
		}).Error
	}); err != nil {
		return
	}
	return x.RefreshCache(ctx)
}
