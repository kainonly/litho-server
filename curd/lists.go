package curd

import (
	"van-api/helper/res"
)

type Lists struct {
	common
	conditions Conditions
	query      Query
	orders     []string
	field      []string
}

func (c *Lists) Where(conditions Conditions) *Lists {
	c.conditions = conditions
	return c
}

func (c *Lists) Query(query Query) *Lists {
	c.query = query
	return c
}

func (c *Lists) OrderBy(orders []string) *Lists {
	c.orders = orders
	return c
}

func (c *Lists) Field(field []string) *Lists {
	c.field = field
	return c
}

func (c *Lists) Result() interface{} {
	body := c.body.(BodyAPI)
	var lists []map[string]interface{}
	tx := c.db.Model(c.model)
	conditions := append(c.conditions, body.GetWhere()...)
	for _, condition := range conditions {
		tx.Where("`"+condition[0].(string)+"` "+condition[1].(string)+" ?", condition[2])
	}
	if c.query != nil {
		c.query(tx)
	}
	orders := append(c.orders, body.GetOrder()...)
	for _, order := range orders {
		tx.Order(order)
	}
	if len(c.field) != 0 {
		tx.Select(c.field)
	}
	page := body.GetPagination()
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
