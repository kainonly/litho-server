package curd

import "van-api/helper/res"

type GetBody struct {
	Id    interface{}
	Where Conditions
	Order []string
}

type getModel struct {
	common
	model      interface{}
	body       GetBody
	conditions Conditions
	query      Query
	orders     []string
	field      []string
}

func (c *getModel) Where(conditions Conditions) *getModel {
	c.conditions = conditions
	return c
}

func (c *getModel) Query(query Query) *getModel {
	c.query = query
	return c
}

func (c *getModel) OrderBy(orders []string) *getModel {
	c.orders = orders
	return c
}

func (c *getModel) Field(field []string) *getModel {
	c.field = field
	return c
}

func (c *getModel) Exec() interface{} {
	data := make(map[string]interface{})
	query := c.db.Model(c.model)
	if c.body.Id != nil {
		query = query.Where("`id` = ?", c.body.Id)
	} else {
		conditions := append(c.conditions, c.body.Where...)
		for _, condition := range conditions {
			query = query.Where("`"+condition[0].(string)+"` "+condition[1].(string)+" ?", condition[2])
		}
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
	query.First(&data)
	return res.Data(data)
}
