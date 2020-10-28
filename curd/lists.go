package curd

import (
	"van-api/helper/res"
)

type ListsBody struct {
	Where Conditions
	Order []string
	Page  Pagination
}

type Pagination struct {
	Index int64
	Limit int64
}

type listsModel struct {
	common
	model      interface{}
	body       ListsBody
	conditions Conditions
	query      Query
	orders     []string
	field      []string
}

func (c *listsModel) Where(conditions Conditions) *listsModel {
	c.conditions = conditions
	return c
}

func (c *listsModel) Query(query Query) *listsModel {
	c.query = query
	return c
}

func (c *listsModel) OrderBy(orders []string) *listsModel {
	c.orders = orders
	return c
}

func (c *listsModel) Field(field []string) *listsModel {
	c.field = field
	return c
}

func (c *listsModel) Exec() interface{} {
	var lists []map[string]interface{}
	query := c.db.Model(c.model)
	conditions := append(c.conditions, c.body.Where...)
	for _, condition := range conditions {
		query = query.Where("`"+condition[0].(string)+"` "+condition[1].(string)+" ?", condition[2])
	}
	if c.query != nil {
		query = c.query(query)
	}
	orders := append(c.orders, c.body.Order...)
	for _, order := range orders {
		query = query.Order(order)
	}
	if len(c.field) != 0 {
		query = query.Select(c.field)
	}
	page := c.body.Page
	if page != (Pagination{}) {
		query = query.
			Limit(int(page.Limit)).
			Offset(int((page.Index - 1) * page.Limit))
	}
	var total int64
	query.Count(&total).Find(&lists)
	return res.Data(map[string]interface{}{
		"lists": lists,
		"total": total,
	})
}
