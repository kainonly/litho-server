package curd

import (
	"van-api/helper/res"
)

type ArrayCondition [][]interface{}

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
	//rType := reflect.SliceOf(reflect.TypeOf(c.model))
	//rSlice := reflect.MakeSlice(rType, 0, 0)
	//rLists := reflect.New(rSlice.Type())
	var lists []map[string]interface{}
	query := c.db.Model(c.model)
	for _, condition := range c.conditions {
		query.Where(condition[0], condition[1])
	}
	for _, order := range c.orders {
		query.Order(order)
	}
	if len(c.field) != 0 {
		query.Select(c.field)
	}
	query.Find(&lists)
	return res.Data(lists)
}
