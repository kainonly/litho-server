package curd

import "gorm.io/gorm"

type common struct {
	db    *gorm.DB
	model interface{}
	body  interface{}
}

type BodyAPI interface {
	GetWhere() [][3]string
	GetOrder() []string
}

type Body struct {
	Where [][3]string
	Order []string
	BodyAPI
}

func (c *Body) GetWhere() [][3]string {
	return c.Where
}

func (c *Body) GetOrder() []string {
	return c.Order
}

type ArrayCondition [][]interface{}
