package roles

import (
	"context"
	"database/sql"
	"litho-api/common"
	"litho-api/model"

	"github.com/bytedance/sonic"
	"github.com/goforj/wire"
	"github.com/kainonly/go/help"
)

var Provides = wire.NewSet(
	wire.Struct(new(Controller), "*"),
	wire.Struct(new(Service), "*"),
)

const Key = "roles"

type Controller struct {
	V *common.Values

	RolesX *Service
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
	return x.RDb.Del(ctx, "iam:roles").Err()
}

func (x *Service) GetIAMRole(ctx context.Context, id string) (result *common.IAMStrategy, err error) {
	var exists int64
	if exists, err = x.RDb.Exists(ctx, "iam:roles").Result(); err != nil {
		return
	}
	if exists != 0 {
		var b []byte
		if b, err = x.RDb.HGet(ctx, "iam:roles", id).Bytes(); err != nil {
			return
		}
		if err = sonic.Unmarshal(b, &result); err != nil {
			return
		}
	} else {
		var rows *sql.Rows
		if rows, err = x.Db.Model(model.Role{}).WithContext(ctx).
			Select(`id`, `strategy`).
			Rows(); err != nil {
			return
		}
		defer rows.Close()
		contents := make(map[string]string)
		for rows.Next() {
			var data model.Role
			if err = x.Db.ScanRows(rows, &data); err != nil {
				return
			}
			if data.ID == id {
				result = &common.IAMStrategy{
					Navs:        data.Strategy["navs"].([]any),
					Routes:      data.Strategy["routes"].([]any),
					Permissions: data.Strategy["permissions"].([]any),
				}
			}
			if contents[data.ID], err = sonic.MarshalString(data.Strategy); err != nil {
				return
			}
		}
		if err = x.RDb.HMSet(ctx, `iam:roles`, contents).Err(); err != nil {
			return
		}
	}
	return
}
