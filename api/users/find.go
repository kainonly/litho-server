package users

import (
	"context"
	"database/sql"
	"litho-api/common"
	"litho-api/model"
	"strconv"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
)

type FindDto struct {
	common.FindDto

	OrgID  string `query:"org_id"`
	RoleID string `query:"role_id"`
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
	ID       string    `json:"id"`
	Status   *bool     `json:"status"`
	RoleIDs  []string  `json:"role_ids" gorm:"-"`
	Roles    []string  `json:"roles" gorm:"-"`
	OrgIDs   []string  `json:"org_ids" gorm:"-"`
	Orgs     []string  `json:"orgs" gorm:"-"`
	Email    string    `json:"email"`
	Phone    string    `json:"phone"`
	Name     string    `json:"name"`
	Avatar   string    `json:"avatar"`
	Sessions int32     `json:"sessions"`
	LastTime time.Time `json:"last_time"`
}

func (x *Service) Find(ctx context.Context, user *common.IAMUser, dto FindDto) (total int64, results []*FindResult, err error) {
	results = make([]*FindResult, 0)
	do := x.Db.Model(model.User{}).WithContext(ctx)

	// 查询条件过滤
	if dto.Q != "" {
		keyword := dto.GetKeyword()
		do = do.Where(
			do.Where(`name like ?`, keyword).
				Or(`phone like ?`, keyword),
		)
	}
	if dto.OrgID != "" {
		do = do.Where(`id in (select user_id from user_org where org_id = ?)`, dto.OrgID)
	}
	if dto.RoleID != "" {
		do = do.Where(`id in (select user_id from user_role where role_id = ?)`, dto.RoleID)
	}

	if err = do.Count(&total).Error; err != nil {
		return
	}

	var rows *sql.Rows
	ctx = common.SetPipe(ctx, common.NewFindPipe().SkipTs().
		Omit("password", "create_time", "update_time"))
	if rows, err = dto.Factory(ctx, do.Order(`id`)).Rows(); err != nil {
		return
	}
	defer rows.Close()

	userIDs := make([]string, 0)
	for rows.Next() {
		var data *FindResult
		if err = x.Db.ScanRows(rows, &data); err != nil {
			return
		}
		userIDs = append(userIDs, data.ID)
		results = append(results, data)
	}

	// 角色关联
	rolesM := make(map[string][]string)
	if len(userIDs) != 0 {
		var roleRows *sql.Rows
		if roleRows, err = x.Db.Model(model.UserRole{}).WithContext(ctx).
			Select(`user_id`, `role_id`).
			Where(`user_id in (?)`, userIDs).
			Rows(); err != nil {
			return
		}
		for roleRows.Next() {
			var data model.UserRole
			if err = x.Db.ScanRows(roleRows, &data); err != nil {
				return
			}
			rolesM[data.UserID] = append(rolesM[data.UserID], data.RoleID)
		}
		roleRows.Close()
	}

	roleIDs := make([]string, 0)
	for _, ids := range rolesM {
		roleIDs = append(roleIDs, ids...)
	}
	var roleM map[string]string
	if roleM, err = x.GetRoleM(ctx, roleIDs); err != nil {
		return
	}

	for _, result := range results {
		result.RoleIDs = rolesM[result.ID]
		result.Roles = make([]string, 0)
		for _, roleId := range result.RoleIDs {
			if v, ok := roleM[roleId]; ok {
				result.Roles = append(result.Roles, v)
			}
		}
	}

	// 组织关联
	orgsM := make(map[string][]string)
	if len(userIDs) != 0 {
		var orgRows *sql.Rows
		if orgRows, err = x.Db.Model(model.UserOrg{}).WithContext(ctx).
			Select(`user_id`, `org_id`).
			Where(`user_id in (?)`, userIDs).
			Order(`user_id, org_id`).
			Rows(); err != nil {
			return
		}
		for orgRows.Next() {
			var data model.UserOrg
			if err = x.Db.ScanRows(orgRows, &data); err != nil {
				return
			}
			orgsM[data.UserID] = append(orgsM[data.UserID], data.OrgID)
		}
		orgRows.Close()
	}

	orgIDs := make([]string, 0)
	for _, ids := range orgsM {
		orgIDs = append(orgIDs, ids...)
	}
	var orgM map[string]*model.Org
	if orgM, err = x.GetOrgM(ctx, orgIDs); err != nil {
		return
	}

	for _, result := range results {
		result.OrgIDs = orgsM[result.ID]
		result.Orgs = make([]string, 0)
		for _, orgId := range result.OrgIDs {
			if v, ok := orgM[orgId]; ok {
				result.Orgs = append(result.Orgs, v.Name)
			}
		}
	}

	return
}
