package curd

import (
	"gorm.io/gorm"
)

type Curd struct {
	db *gorm.DB
}

func Initialize(db *gorm.DB) *Curd {
	c := new(Curd)
	c.db = db
	return c
}

func (c *Curd) Originlists(model interface{}, body interface{}) *OriginLists {
	m := new(OriginLists)
	m.db = c.db
	m.model = model
	m.body = body
	return m
}

func (c *Curd) Get(model interface{}, body interface{}) *Get {
	m := new(Get)
	m.db = c.db
	m.model = model
	m.body = body
	return m
}
