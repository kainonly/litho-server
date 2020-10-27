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

func (c *edit) Result() interface{} {
	return res.Ok()
}
