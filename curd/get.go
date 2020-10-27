package curd

import "van-api/helper/res"

type GetBody struct {
	Id    interface{}
	Where Conditions
	Order []string
}

type get struct {
	common
	model      interface{}
	body       GetBody
	conditions Conditions
	query      Query
	orders     []string
	field      []string
}

func (c *get) Where(conditions Conditions) *get {
	c.conditions = conditions
	return c
}

func (c *get) Query(query Query) *get {
	c.query = query
	return c
}

func (c *get) OrderBy(orders []string) *get {
	c.orders = orders
	return c
}

func (c *get) Field(field []string) *get {
	c.field = field
	return c
}

func (c *get) Result() interface{} {
	data := make(map[string]interface{})
	tx := c.db.Model(c.model)
	if c.body.Id != nil {
		tx.Where("`id` = ?", c.body.Id)
	} else {
		conditions := append(c.conditions, c.body.Where...)
		for _, condition := range conditions {
			tx.Where("`"+condition[0].(string)+"` "+condition[1].(string)+" ?", condition[2])
		}
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
	tx.First(&data)
	return res.Data(data)
}
