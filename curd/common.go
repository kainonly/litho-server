package curd

import "gorm.io/gorm"

type common struct {
	db    *gorm.DB
	model interface{}
	body  interface{}
}

type ArrayCondition [][]interface{}

type BodyAPI interface {
	GetWhere() ArrayCondition
	GetOrder() []string
}

type Body struct {
	Where ArrayCondition
	Order []string
	BodyAPI
}

func (c *Body) GetWhere() ArrayCondition {
	return c.Where
}

func (c *Body) GetOrder() []string {
	return c.Order
}
