package curd

import (
	"gorm.io/gorm"
	"van-api/helper/res"
)

type add struct {
	common
	model interface{}
	after func(tx *gorm.DB) error
}

func (c *add) After(hook func(tx *gorm.DB) error) *add {
	c.after = hook
	return c
}

func (c *add) Result() interface{} {
	if c.after == nil {
		if err := c.db.Create(c.model).Error; err != nil {
			return err
		}
	} else {
		err := c.db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Create(c.model).Error; err != nil {
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
