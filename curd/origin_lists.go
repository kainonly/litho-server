package curd

import (
	"van-api/helper/res"
)

type OriginLists struct {
	common
	conditions ArrayCondition
	orders     []string
	field      []string
}

func (c *OriginLists) Where(conditions ArrayCondition) *OriginLists {
	c.conditions = conditions
	return c
}

func (c *OriginLists) OrderBy(orders []string) *OriginLists {
	c.orders = orders
	return c
}

func (c *OriginLists) Field(field []string) *OriginLists {
	c.field = field
	return c
}

func (c *OriginLists) Result() interface{} {
	var lists []map[string]interface{}
	query := c.db.Model(c.model)
	conditions := append(c.conditions, c.body.(BodyAPI).GetWhere()...)
	for _, condition := range conditions {
		query.Where("`"+condition[0].(string)+"` "+condition[1].(string)+" ?", condition[2])
	}
	orders := append(c.orders, c.body.(BodyAPI).GetOrder()...)
	for _, order := range orders {
		query.Order(order)
	}
	if len(c.field) != 0 {
		query.Select(c.field)
	}
	query.Find(&lists)
	return res.Data(lists)
}
