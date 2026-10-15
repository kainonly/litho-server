package routes

import (
	"context"
	"database/sql"
	"litho-api/common"
	"litho-api/model"
	"strings"

	"github.com/bytedance/sonic"
	"github.com/goforj/wire"
	"github.com/kainonly/go/help"
)

var Provides = wire.NewSet(
	wire.Struct(new(Controller), "*"),
	wire.Struct(new(Service), "*"),
)

const Key = "routes"

type Controller struct {
	V *common.Values

	RoutesX *Service
}

type Service struct {
	*common.Inject
}

type M = map[string]any

var ErrAccess = help.E(0, `您不具备系统权限，请联系超级管理员进行管理`)

func (x *Service) CheckAccess(ctx context.Context, user *common.IAMUser, ids ...string) (err error) {
	if err = user.Can(`WM`); err != nil {
		return ErrAccess
	}
	return
}

func (x *Service) RefreshCache(ctx context.Context) error {
	return x.RDb.Del(ctx, "iam:routes").Err()
}

func (x *Service) GetIAMRoute(ctx context.Context, link string) (result *common.IAMRoute, err error) {
	var exists int64
	if exists, err = x.RDb.Exists(ctx, "iam:routes").Result(); err != nil {
		return
	}
	field := strings.Join(strings.Split(link, "/")[:4], "/") // 仅验证前端路由前三层
	if exists != 0 {
		var b []byte
		if b, err = x.RDb.HGet(ctx, "iam:routes", field).Bytes(); err != nil {
			return
		}
		if err = sonic.Unmarshal(b, &result); err != nil {
			return
		}
	} else {
		var rows *sql.Rows
		if rows, err = x.Db.Model(model.Route{}).WithContext(ctx).
			Select(`id`, `link`, `strategy`).
			Rows(); err != nil {
			return
		}
		defer rows.Close()
		contents := make(map[string]string)
		for rows.Next() {
			var data *common.IAMRoute
			if err = x.Db.ScanRows(rows, &data); err != nil {
				return
			}
			if data.Link == field {
				result = data
			}
			if contents[data.Link], err = sonic.MarshalString(data); err != nil {
				return
			}
		}
		if err = x.RDb.HMSet(ctx, `iam:routes`, contents).Err(); err != nil {
			return
		}
	}
	return
}
