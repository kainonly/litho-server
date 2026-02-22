package index

import (
	"database/sql"
	"server/common"
	"server/model"

	"context"

	"github.com/cloudwego/hertz/pkg/app"
)

func (x *Controller) GetLayout(ctx context.Context, c *app.RequestContext) {
	user := common.GetIAM(c)
	result, err := x.IndexX.GetLayout(ctx, user)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, result)
}

type LayoutResult struct {
	Navs     []string           `json:"navs"`
	NavMenus map[string][]*Menu `json:"nav_menus"`
}

type Nav struct {
	Key  string `json:"key"`
	Name string `json:"name"`
	Icon string `json:"icon"`
	Link string `json:"link"`
}

type Menu struct {
	ID       string  `json:"-"`
	Pid      string  `json:"-"`
	Name     string  `json:"name"`
	Disabled bool    `json:"disabled"`
	Icon     string  `json:"icon,omitempty"`
	Link     string  `json:"link,omitempty"`
	Children []*Menu `json:"children,omitempty"`
}

func (x *Service) GetLayout(ctx context.Context, user *common.IAMUser) (result *LayoutResult, err error) {
	result = &LayoutResult{
		Navs:     user.Strategy.Navs,
		NavMenus: make(map[string][]*Menu),
	}
	// 获取授权子路由
	if result.NavMenus, err = x.GetLayoutNavMenus(ctx, user.Strategy); err != nil {
		return
	}
	return
}

func (x *Service) GetLayoutNavMenus(ctx context.Context, strategy *common.RoleStrategy) (navMenus map[string][]*Menu, err error) {
	navMenus = make(map[string][]*Menu)
	var rows *sql.Rows
	if rows, err = x.Db.Model(model.Route{}).WithContext(ctx).
		Select(`id`, `active`, `nav`, `pid`, `name`, `icon`, `link`).
		Where(`id in (?)`, strategy.Routes).
		Order(`type`).
		Order(`sort`).
		Rows(); err != nil {
		return
	}
	defer rows.Close()

	navDict := make(map[string][]*Menu)
	nodeM := make(map[string]*Menu)
	for rows.Next() {
		var data model.Route
		if err = x.Db.ScanRows(rows, &data); err != nil {
			return
		}
		if navDict[data.Nav] == nil {
			navDict[data.Nav] = make([]*Menu, 0)
		}
		if data.Pid == "0" {
			nodeM[data.ID] = &Menu{
				ID:       data.ID,
				Pid:      data.Pid,
				Name:     data.Name,
				Disabled: !*data.Active,
				Children: make([]*Menu, 0),
			}
			navDict[data.Nav] = append(navDict[data.Nav], nodeM[data.ID])
		} else {
			parent := nodeM[data.Pid]
			if parent != nil {
				parent.Children = append(parent.Children, &Menu{
					ID:       data.ID,
					Pid:      data.Pid,
					Name:     data.Name,
					Disabled: !*data.Active,
					Icon:     data.Icon,
					Link:     data.Link,
				})
			}
		}
	}
	for nav, menus := range navDict {
		navMenus[nav] = menus
	}
	return
}
