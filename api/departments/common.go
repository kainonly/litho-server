package departments

import (
	"context"
	"database/sql"
	"server/common"
	"server/model"

	"github.com/goforj/wire"
)

const (
	Key   = "departments"
	Label = "部门"
)

var Provides = wire.NewSet(
	wire.Struct(new(Controller), "*"),
	wire.Struct(new(Service), "*"),
)

type Controller struct {
	DepartmentsX *Service
}

type Service struct {
	*common.Inject
}

func (x *Service) GetDepartmentM(ctx context.Context, ids []string) (result map[string]*model.Department, err error) {
	result = make(map[string]*model.Department)
	var rows *sql.Rows
	if rows, err = x.Db.Model(model.Department{}).WithContext(ctx).
		Select([]string{"id", "type", "name"}).
		Where(`id in (?)`, ids).
		Rows(); err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var data *model.Department
		if err = x.Db.ScanRows(rows, &data); err != nil {
			return
		}
		result[data.ID] = data
	}
	return
}
