package users

import (
	"context"
	"database/sql"
	"server/common"
	"server/model"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
)

type FindDto struct {
	common.FindDto
}

func (x *Controller) Find(ctx context.Context, c *app.RequestContext) {
	var dto FindDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	user := common.GetIAM(c)
	total, data, err := x.UsersX.Find(ctx, user, dto)
	if err != nil {
		c.Error(err)
		return
	}

	c.Header("x-total", strconv.Itoa(int(total)))
	c.JSON(200, data)
}

type FindResult struct {
	ID       string `json:"id"`
	Status   *bool  `json:"status"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Sessions int32  `json:"sessions"`
}

func (x *Service) Find(ctx context.Context, user *common.IAMUser, dto FindDto) (total int64, results []*FindResult, err error) {
	results = make([]*FindResult, 0)
	do := x.Db.Model(model.User{}).WithContext(ctx)

	if dto.Q != "" {
	}

	if err = do.Count(&total).Error; err != nil {
		return
	}

	var rows *sql.Rows
	ctx = common.SetPipe(ctx, common.NewFindPipe().SkipTs().
		Omit("password", "create_time", "update_time"))
	if rows, err = dto.Factory(ctx, do).Rows(); err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var data *FindResult
		if err = x.Db.ScanRows(rows, &data); err != nil {
			return
		}
		results = append(results, data)
	}

	return
}
