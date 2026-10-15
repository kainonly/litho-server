package users

import (
	"context"
	"fmt"
	"litho-api/common"
	"litho-api/model"

	"github.com/cloudwego/hertz/pkg/app"
)

type SearchDto struct {
	common.SearchDto

	OrgID string `query:"org_id"`
}

func (x *Controller) Search(ctx context.Context, c *app.RequestContext) {
	var dto SearchDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	user := common.GetIAM(c)
	data, err := x.UsersX.Search(ctx, user, dto)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, data)
}

type SearchResult struct {
	ID     string   `json:"id"`
	OrgIDs []string `json:"org_ids" gorm:"-"`
	Orgs   []string `json:"orgs" gorm:"-"`
	Name   string   `json:"name"`
	Phone  string   `json:"phone"`
}

func (x *Service) Search(ctx context.Context, user *common.IAMUser, dto SearchDto) (results []SearchResult, err error) {
	do := x.Db.Model(model.User{}).WithContext(ctx)

	if dto.OrgID != "" {
		do = do.Where(`id in (select user_id from user_org where org_id = ?)`, dto.OrgID)
	}
	if dto.Q != "" {
		do = do.Where(`name like ?`, fmt.Sprintf(`%%%s%%`, dto.Q))
	}

	results = make([]SearchResult, 0, len(results))
	ctx = common.SetPipe(ctx, common.NewSearchPipe(`id`, `name`, `phone`).SkipAsync())
	if err = dto.Find(ctx, do, &results); err != nil {
		return
	}

	userIDs := make([]string, 0, len(results))
	for i := range results {
		results[i].OrgIDs = make([]string, 0)
		results[i].Orgs = make([]string, 0)
		userIDs = append(userIDs, results[i].ID)
	}
	if len(userIDs) == 0 {
		return
	}
	var relations []model.UserOrg
	if err = x.Db.Model(model.UserOrg{}).WithContext(ctx).
		Where(`user_id in (?)`, userIDs).
		Order(`user_id, org_id`).
		Find(&relations).Error; err != nil {
		return
	}
	orgIDs := make([]string, 0, len(relations))
	orgsM := make(map[string][]string)
	for _, relation := range relations {
		orgsM[relation.UserID] = append(orgsM[relation.UserID], relation.OrgID)
		orgIDs = append(orgIDs, relation.OrgID)
	}
	orgM := make(map[string]*model.Org)
	if len(orgIDs) != 0 {
		if orgM, err = x.GetOrgM(ctx, orgIDs); err != nil {
			return nil, err
		}
	}
	for i := range results {
		if ids, ok := orgsM[results[i].ID]; ok {
			results[i].OrgIDs = ids
		}
		for _, orgID := range results[i].OrgIDs {
			if org, ok := orgM[orgID]; ok {
				results[i].Orgs = append(results[i].Orgs, org.Name)
			}
		}
	}
	return
}
