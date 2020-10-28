package curd

import (
	"gorm.io/gorm"
	"van-api/helper/res"
)

type addModel struct {
	common
	after func(tx *gorm.DB) error
}

func (c *addModel) After(hook func(tx *gorm.DB) error) *addModel {
	c.after = hook
	return c
}

func (c *addModel) Exec(value interface{}) interface{} {
	query := c.db
	if c.after == nil {
		if err := query.Create(&value).Error; err != nil {
			return err
		}
	} else {
		err := query.Transaction(func(tx *gorm.DB) error {
			if err := tx.Create(&value).Error; err != nil {
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
