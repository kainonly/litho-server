package curd

import (
	"gorm.io/gorm"
	"van-api/helper/res"
)

type DeleteBody struct {
	Id    interface{}
	Where Conditions
}

type deleteModel struct {
	common
	model      interface{}
	body       DeleteBody
	conditions Conditions
	query      Query
	prep       func(tx *gorm.DB) error
	after      func(tx *gorm.DB) error
}

func (c *deleteModel) Where(conditions Conditions) *deleteModel {
	c.conditions = conditions
	return c
}

func (c *deleteModel) Query(query Query) *deleteModel {
	c.query = query
	return c
}

func (c *deleteModel) Prep(hook func(tx *gorm.DB) error) *deleteModel {
	c.prep = hook
	return c
}

func (c *deleteModel) After(hook func(tx *gorm.DB) error) *deleteModel {
	c.after = hook
	return c
}

func (c *deleteModel) Exec() interface{} {
	query := c.db
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
	if c.after == nil && c.prep == nil {
		if err := query.Delete(c.model).Error; err != nil {
			return err
		}
	} else {
		err := query.Transaction(func(tx *gorm.DB) error {
			if err := c.prep(tx); err != nil {
				return err
			}
			if err := tx.Delete(c.model).Error; err != nil {
				return err
			}
			if err := c.after(tx); err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return res.Error(err)
		}
	}
	return res.Ok()
}
