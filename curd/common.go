package curd

import "gorm.io/gorm"

type common struct {
	db    *gorm.DB
	model interface{}
	body  interface{}
}

type Conditions [][]interface{}
type Query func(tx *gorm.DB)

type BodyAPI interface {
	GetWhere() Conditions
	GetOrder() []string
}

type Body struct {
	Where Conditions
	Order []string
	BodyAPI
}

func (c *Body) GetWhere() Conditions {
	return c.Where
}

func (c *Body) GetOrder() []string {
	return c.Order
}
