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

type lists struct {
	common
	model      interface{}
	body       ListsBody
	conditions Conditions
	query      Query
	orders     []string
	field      []string
}

func (c *lists) Where(conditions Conditions) *lists {
	c.conditions = conditions
	return c
}

func (c *lists) Query(query Query) *lists {
	c.query = query
	return c
}

func (c *lists) OrderBy(orders []string) *lists {
	c.orders = orders
	return c
}

func (c *lists) Field(field []string) *lists {
	c.field = field
	return c
}

func (c *lists) Result() interface{} {
	var lists []map[string]interface{}
	tx := c.db.Model(c.model)
	conditions := append(c.conditions, c.body.Where...)
	for _, condition := range conditions {
		tx.Where("`"+condition[0].(string)+"` "+condition[1].(string)+" ?", condition[2])
	}
	if c.query != nil {
		c.query(tx)
	}
	orders := append(c.orders, c.body.Order...)
	for _, order := range orders {
		tx.Order(order)
	}
	if len(c.field) != 0 {
		tx.Select(c.field)
	}
	page := c.body.Page
	if page != (Pagination{}) {
		tx.Limit(int(page.Limit))
		tx.Offset(int((page.Index - 1) * page.Limit))
	}
	var total int64
	tx.Count(&total)
	tx.Find(&lists)
	return res.Data(map[string]interface{}{
		"lists": lists,
		"total": total,
	})
}
