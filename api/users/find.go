package users

import (
	"context"
	"database/sql"
	"strconv"

	"server/common"
	"server/model"

	"github.com/cloudwego/hertz/pkg/app"
)

type FindDto struct {
	common.FindDto

	RoleID string `json:"role_id"`
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
	ID     string `json:"id"`
	Active bool   `json:"active"`
	Email  string `json:"email"`
	Phone  string `json:"phone"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

func (x *Service) Find(ctx context.Context, user *common.IAMUser, dto FindDto) (total int64, results []*FindResult, err error) {

	do := x.Db.Model(model.User{}).WithContext(ctx)
	if dto.RoleID != "" {
	}

	if dto.Q != "" {
		keyword := dto.GetKeyword()
		do = do.Where(
			do.Where(`email like ?`, keyword).
				Or(`name like ?`, keyword),
		)
	}

	if err = do.Count(&total).Error; err != nil {
		return
	}

	var rows *sql.Rows
	ctx = common.SetPipe(ctx, common.NewFindPipe().SkipTs().
		Omit(`created_at`, `updated_at`, `password`))
	db, err := dto.Factory(ctx, do)
	if err != nil {
		return
	}
	if rows, err = db.Rows(); err != nil {
		return
	}
	defer rows.Close()

	results = make([]*FindResult, 0)
	for rows.Next() {
		var data *FindResult
		if err = x.Db.ScanRows(rows, &data); err != nil {
			return
		}
		results = append(results, data)
	}

	return
}
