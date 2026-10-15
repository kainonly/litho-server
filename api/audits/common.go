package audits

import (
	"context"
	"litho-api/common"
	"litho-api/model"

	"github.com/goforj/wire"
)

var Provides = wire.NewSet(
	wire.Struct(new(Controller), "*"),
	wire.Struct(new(Service), "*"),
)

type Controller struct {
	AuditsX *Service
}

type Service struct {
	*common.Inject
}

type Enum struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

type UserInfo struct {
	ID     string
	Name   string
	Phone  string
	OrgIDs []string
	Orgs   []string
}

type userOrgRow struct {
	UserID string
	OrgID  string
	Name   string
}

func (x *Service) GetUserM(ctx context.Context, ids []string) (map[string]*UserInfo, error) {
	result := make(map[string]*UserInfo)
	if len(ids) == 0 {
		return result, nil
	}
	var users []UserInfo
	if err := x.Db.Model(model.User{}).WithContext(ctx).
		Select([]string{"id", "name", "phone"}).
		Where("id in (?)", ids).
		Find(&users).Error; err != nil {
		return nil, err
	}
	for i := range users {
		users[i].OrgIDs = make([]string, 0)
		users[i].Orgs = make([]string, 0)
		result[users[i].ID] = &users[i]
	}

	var orgRows []userOrgRow
	if err := x.Db.WithContext(ctx).
		Table(model.TableNameUserOrg+" as uo").
		Select("uo.user_id, org.id as org_id, org.name").
		Joins("join "+model.TableNameOrg+" as org on org.id = uo.org_id").
		Where("uo.user_id in (?)", ids).
		Order("uo.user_id, org.id").
		Scan(&orgRows).Error; err != nil {
		return nil, err
	}
	for _, row := range orgRows {
		if user, ok := result[row.UserID]; ok {
			user.OrgIDs = append(user.OrgIDs, row.OrgID)
			user.Orgs = append(user.Orgs, row.Name)
		}
	}
	return result, nil
}
