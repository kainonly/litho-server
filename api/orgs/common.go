package orgs

import (
	"context"
	"database/sql"
	"server/common"
	"server/model"

	"github.com/goforj/wire"
)

const (
	Key = "orgs"
	Label    = "组织"
)

var Provides = wire.NewSet(
	wire.Struct(new(Controller), "*"),
	wire.Struct(new(Service), "*"),
)

type Controller struct {
	OrgsX *Service
}

type Service struct {
	*common.Inject
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
