package users

import (
	"context"
	"litho-api/common"
	"litho-api/model"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/kainonly/go/help"
	"gorm.io/gorm"
)

type SetOrgsDto struct {
	ID     string   `json:"id" vd:"required"`
	OrgIDs []string `json:"org_ids" vd:"required"`
}

func (x *Controller) SetOrgs(ctx context.Context, c *app.RequestContext) {
	var dto SetOrgsDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	user := common.GetIAM(c)
	if err := x.UsersX.SetOrgs(ctx, user, dto); err != nil {
		c.Error(err)
		return
	}

	common.SetTracking(c, Key, "set_orgs", dto.ID)
	c.JSON(200, help.Ok())
}

func (x *Service) SetOrgs(ctx context.Context, user *common.IAMUser, dto SetOrgsDto) (err error) {
	if err = x.CheckAccess(ctx, user); err != nil {
		return
	}
	if len(dto.OrgIDs) == 0 {
		return help.E(0, "至少选择一个组织")
	}

	seen := make(map[string]bool, len(dto.OrgIDs))
	relations := make([]model.UserOrg, 0, len(dto.OrgIDs))
	for _, orgID := range dto.OrgIDs {
		if orgID == "" || seen[orgID] {
			return help.E(0, "组织选择无效或重复")
		}
		seen[orgID] = true
		relations = append(relations, model.UserOrg{UserID: dto.ID, OrgID: orgID})
	}

	if err = x.Db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var target model.User
		if err := tx.Select("id").Where("id = ?", dto.ID).Take(&target).Error; err != nil {
			return err
		}
		var count int64
		if err := tx.Model(&model.Org{}).Where("id IN ?", dto.OrgIDs).Count(&count).Error; err != nil {
			return err
		}
		if count != int64(len(dto.OrgIDs)) {
			return help.E(0, "所选组织不存在")
		}
		if err := tx.Where("user_id = ?", dto.ID).Delete(&model.UserOrg{}).Error; err != nil {
			return err
		}
		if err := tx.Create(&relations).Error; err != nil {
			return err
		}
		return tx.Model(&model.User{}).Where("id = ?", dto.ID).Update("update_time", time.Now()).Error
	}); err != nil {
		return
	}
	return x.RefreshCache(ctx)
}
