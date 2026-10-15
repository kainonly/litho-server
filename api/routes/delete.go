package routes

import (
	"context"
	"database/sql"
	"fmt"
	"litho-api/common"
	"litho-api/model"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/kainonly/go/help"
	"gorm.io/gorm"
)

func (x *Controller) Delete(ctx context.Context, c *app.RequestContext) {
	var dto common.DeleteDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	user := common.GetIAM(c)
	if err := x.RoutesX.Delete(ctx, user, dto); err != nil {
		c.Error(err)
		return
	}

	common.SetTracking(c, Key, "delete", dto.IDs...)
	c.JSON(200, help.Ok())
}

func (x *Service) Delete(ctx context.Context, user *common.IAMUser, dto common.DeleteDto) (err error) {
	if err = x.CheckAccess(ctx, user); err != nil {
		return
	}

	deleteIDM := make(map[string]bool)
	for _, id := range dto.IDs {
		deleteIDM[id] = true
	}
	var rows *sql.Rows
	if rows, err = x.Db.Model(model.Role{}).WithContext(ctx).
		Select([]string{`id`, `strategy`}).
		Where(fmt.Sprintf(`strategy->'routes' ?| array['%s']`, strings.Join(dto.IDs, `','`))).
		Rows(); err != nil {
		return
	}
	strategyM := make(map[string][]string)
	for rows.Next() {
		var role *model.Role
		if err = x.Db.ScanRows(rows, &role); err != nil {
			return
		}
		strategyM[role.ID] = make([]string, 0)
		for _, v := range role.Strategy["routes"].([]any) {
			id := v.(string)
			if !deleteIDM[id] {
				strategyM[role.ID] = append(strategyM[role.ID], id)
			}
		}
	}
	rows.Close()
	return x.Db.Transaction(func(tx *gorm.DB) (errX error) {
		if err = tx.WithContext(ctx).
			Delete(model.Route{}, dto.IDs).Error; err != nil {
			return
		}
		for id, routes := range strategyM {
			if err = tx.Model(model.Role{}).WithContext(ctx).
				Where(`id = ?`, id).
				Update(`strategy`,
					gorm.Expr(`jsonb_set("strategy", array ['routes'], ?)`,
						fmt.Sprintf(`["%s"]`, strings.Join(routes, `","`)),
					),
				).
				Error; err != nil {
				return
			}
		}
		// 权限组变动，刷新缓存
		if err = x.RDb.Del(ctx, "iam:roles").Err(); err != nil {
			return
		}
		return x.RefreshCache(ctx)
	})
}
