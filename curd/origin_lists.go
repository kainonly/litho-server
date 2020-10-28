package curd

import (
	"van-api/helper/res"
)

type OriginListsBody struct {
	Where Conditions
	Order []string
}

type originListsModel struct {
	common
	model      interface{}
	body       OriginListsBody
	conditions Conditions
	query      Query
	orders     []string
	field      []string
}

func (c *originListsModel) Where(conditions Conditions) *originListsModel {
	c.conditions = conditions
	return c
}

func (c *originListsModel) Query(query Query) *originListsModel {
	c.query = query
	return c
}

func (c *originListsModel) OrderBy(orders []string) *originListsModel {
	c.orders = orders
	return c
}

func (c *originListsModel) Field(field []string) *originListsModel {
	c.field = field
	return c
}

func (c *originListsModel) Exec() interface{} {
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
	query.Find(&lists)
	return res.Data(lists)
}
