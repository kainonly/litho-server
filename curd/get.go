package curd

import "van-api/helper/res"

type Get struct {
	common
	conditions Conditions
	query      Query
	orders     []string
	field      []string
}

func (c *Get) Where(conditions Conditions) *Get {
	c.conditions = conditions
	return c
}

func (c *Get) Query(query Query) *Get {
	c.query = query
	return c
}

func (c *Get) OrderBy(orders []string) *Get {
	c.orders = orders
	return c
}

func (c *Get) Field(field []string) *Get {
	c.field = field
	return c
}

func (c *Get) Result() interface{} {
	body := c.body.(BodyAPI)
	data := make(map[string]interface{})
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
	tx.First(&data)
	return res.Data(data)
}
