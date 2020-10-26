package curd

import "gorm.io/gorm"

type common struct {
	db    *gorm.DB
	model interface{}
	body  interface{}
}

func (c *common) initialize(db *gorm.DB, model interface{}, body interface{}) {
	c.db = db
	c.model = model
	c.body = body
}

type Conditions [][]interface{}
type Query func(tx *gorm.DB)

type BodyAPI interface {
	GetWhere() Conditions
	GetOrder() []string
	GetPagination() Pagination
}

type Body struct {
	Where Conditions
	Order []string
	Page  Pagination
	BodyAPI
}

type Pagination struct {
	Index int64
	Limit int64
}

func (c *Body) GetWhere() Conditions {
	return c.Where
}

func (c *Body) GetOrder() []string {
	return c.Order
}

func (c *Body) GetPagination() Pagination {
	return c.Page
}
