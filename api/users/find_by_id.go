package users

import (
	"context"
	"litho-api/common"
	"litho-api/model"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
)

func (x *Controller) FindById(ctx context.Context, c *app.RequestContext) {
	var dto common.FindByIdDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	user := common.GetIAM(c)
	data, err := x.UsersX.FindById(ctx, user, dto)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, data)
}

type FindByIdResult struct {
	ID       string    `json:"id"`
	Status   *bool     `json:"status"`
	OrgID    string    `json:"org_id"`
	OrgIDs   []string  `json:"org_ids" gorm:"-"`
	Name     string    `json:"name"`
	Phone    string    `json:"phone"`
	Email    string    `json:"email"`
	Avatar   string    `json:"avatar"`
	Sessions int32     `json:"sessions"`
	LastTime time.Time `json:"last_time"`
}

func (x *Service) FindById(ctx context.Context, user *common.IAMUser, dto common.FindByIdDto) (result FindByIdResult, err error) {
	do := x.Db.Model(model.User{}).WithContext(ctx)
	ctx = common.SetPipe(ctx, common.NewFindByIdPipe().SkipTs().
		Omit(`create_time`, `update_time`, `password`))
	if err = dto.Take(ctx, do, &result); err != nil {
		return
	}

	// 组织关联（用户可属于多个组织；保留 org_id 兼容旧调用方）
	orgIDs := make([]string, 0)
	if err = x.Db.Model(model.UserOrg{}).WithContext(ctx).
		Where(`user_id = ?`, result.ID).
		Pluck(`org_id`, &orgIDs).Error; err != nil {
		return
	}
	if len(orgIDs) != 0 {
		result.OrgID = orgIDs[0]
	}
	result.OrgIDs = orgIDs
	return
}
