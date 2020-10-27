package curd

import (
	"gorm.io/gorm"
	"van-api/helper/res"
)

type EditBody struct {
	Id    interface{}
	Where Conditions
}

type edit struct {
	common
	model      interface{}
	body       EditBody
	conditions Conditions
	query      Query
	after      func(tx *gorm.DB) error
}

func (c *edit) Where(conditions Conditions) *edit {
	c.conditions = conditions
	return c
}

func (c *edit) Query(query Query) *edit {
	c.query = query
	return c
}

func (c *edit) After(hook func(tx *gorm.DB) error) *edit {
	c.after = hook
	return c
}

func (c *edit) Result(update interface{}) interface{} {
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
	if c.after == nil {
		if err := tx.Updates(update); err != nil {
			return err
		}
	} else {
		err := c.db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Updates(update).Error; err != nil {
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
