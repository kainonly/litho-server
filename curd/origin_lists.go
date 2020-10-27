package curd

import (
	"van-api/helper/res"
)

type OriginListsBody struct {
	Where Conditions
	Order []string
}

type originLists struct {
	common
	model      interface{}
	body       OriginListsBody
	conditions Conditions
	query      Query
	orders     []string
	field      []string
}

func (c *originLists) Where(conditions Conditions) *originLists {
	c.conditions = conditions
	return c
}

func (c *originLists) Query(query Query) *originLists {
	c.query = query
	return c
}

func (c *originLists) OrderBy(orders []string) *originLists {
	c.orders = orders
	return c
}

func (c *originLists) Field(field []string) *originLists {
	c.field = field
	return c
}

func (c *originLists) Result() interface{} {
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
	tx.Find(&lists)
	return res.Data(lists)
}
