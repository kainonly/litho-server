package users

import (
	"context"
	"database/sql"
	"litho-api/api/sessions"
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

const Key = "users"

type Controller struct {
	V *common.Values

	UsersX *Service
}

type Service struct {
	*common.Inject

	SessionsX *sessions.Service
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
	return x.RDb.Del(ctx, "iam:users").Err()
}

func (x *Service) GetIAMUser(ctx context.Context, id string) (result *common.IAMUser, err error) {
	var exists int64
	if exists, err = x.RDb.Exists(ctx, "iam:users").Result(); err != nil {
		return
	}

	if exists != 0 {
		var b []byte
		if b, err = x.RDb.HGet(ctx, "iam:users", id).Bytes(); err != nil {
			return
		}
		if err = sonic.Unmarshal(b, &result); err != nil {
			return
		}
	} else {
		var rows *sql.Rows
		if rows, err = x.Db.Model(model.User{}).WithContext(ctx).
			Select(`id`, `status`).
			Rows(); err != nil {
			return
		}
		defer rows.Close()
		users := make([]*common.IAMUser, 0)
		for rows.Next() {
			var user *common.IAMUser
			if err = x.Db.ScanRows(rows, &user); err != nil {
				return
			}
			user.RoleID = "0"
			users = append(users, user)
		}

		// 角色关联
		rolesM := make(map[string]string)
		if rows, err = x.Db.Model(model.UserRole{}).WithContext(ctx).
			Select(`user_id`, `role_id`).
			Rows(); err != nil {
			return
		}
		for rows.Next() {
			var data model.UserRole
			if err = x.Db.ScanRows(rows, &data); err != nil {
				return
			}
			rolesM[data.UserID] = data.RoleID
		}
		rows.Close()

		contents := make(map[string]string)
		for _, user := range users {
			if roleID, ok := rolesM[user.ID]; ok {
				user.RoleID = roleID
			}
			if user.ID == id {
				result = user
			}
			if contents[user.ID], err = sonic.MarshalString(user); err != nil {
				return
			}
		}
		if err = x.RDb.HMSet(ctx, `iam:users`, contents).Err(); err != nil {
			return
		}
	}
	return
}

func (x *Service) GetOrgM(ctx context.Context, ids []string) (result map[string]*model.Org, err error) {
	result = make(map[string]*model.Org)
	var rows *sql.Rows
	if rows, err = x.Db.Model(model.Org{}).WithContext(ctx).
		Select([]string{"id", "type", "name"}).
		Where(`id in (?)`, ids).
		Rows(); err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var data *model.Org
		if err = x.Db.ScanRows(rows, &data); err != nil {
			return
		}
		result[data.ID] = data
	}
	return
}

func (x *Service) GetRoleM(ctx context.Context, ids []string) (result map[string]string, err error) {
	result = make(map[string]string)
	var rows *sql.Rows
	if rows, err = x.Db.Model(model.Role{}).WithContext(ctx).
		Select([]string{"id", "name"}).
		Where(`id in (?)`, ids).
		Rows(); err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var data model.Role
		if err = x.Db.ScanRows(rows, &data); err != nil {
			return
		}
		result[data.ID] = data.Name
	}

	return
}
