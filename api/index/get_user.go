package index

import (
	"context"
	"server/common"
	"server/model"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
)

func (x *Controller) GetUser(ctx context.Context, c *app.RequestContext) {
	user := common.GetIAM(c)
	data, err := x.IndexX.GetUser(ctx, user.ID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, data)
}

type UserResult struct {
	ID         string    `json:"id"`
	Email      string    `json:"email"`
	Sessions   int32     `json:"sessions"`
	CreateTime time.Time `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`
}

func (x *Service) GetUser(ctx context.Context, userId string) (result *UserResult, err error) {
	var data *model.User
	if err = x.Db.Model(model.User{}).WithContext(ctx).
		Where(`id = ?`, userId).
		Take(&data).Error; err != nil {
		return
	}

	result = &UserResult{
		ID:         data.ID,
		Email:      data.Email,
		Sessions:   data.Sessions,
		CreateTime: data.CreateTime,
		UpdateTime: data.UpdateTime,
	}
	return
}
