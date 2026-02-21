package users

import (
    "context"
    "database/sql"
    "server/api/orgs"
    "server/api/sessions"
    "server/common"
    "server/model"

    "github.com/bytedance/sonic"
    "github.com/goforj/wire"
)

const (
    Key = "users"
    Label    = "团队成员"
)

var Provides = wire.NewSet(
    wire.Struct(new(Controller), "*"),
    wire.Struct(new(Service), "*"),
)

type Controller struct {
    UsersX *Service
}

type Service struct {
    *common.Inject

    SessionsX *sessions.Service
    OrgsX     *orgs.Service
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
            Select(`id`, `role_id`, `org_id`, `active`).
            Rows(); err != nil {
            return
        }
        defer rows.Close()
        users := make([]*common.IAMUser, 0)
        orgIDs := make([]string, 0)
        for rows.Next() {
            var user *common.IAMUser
            if err = x.Db.ScanRows(rows, &user); err != nil {
                return
            }
            orgIDs = append(orgIDs, user.OrgID)
            users = append(users, user)
        }

        var orgM map[string]*model.Org
        if orgM, err = x.OrgsX.GetOrgM(ctx, orgIDs); err != nil {
            return
        }

        contents := make(map[string]string)
        for _, user := range users {
            if v, ok := orgM[user.OrgID]; ok {
                user.OrgType = *v.Type
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
