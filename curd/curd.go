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
	m.initialize(c.db, model, body)
	return m
}

func (c *Curd) Lists(model interface{}, body interface{}) *Lists {
	m := new(Lists)
	m.initialize(c.db, model, body)
	return m
}

func (c *Curd) Get(model interface{}, body interface{}) *Get {
	m := new(Get)
	m.initialize(c.db, model, body)
	return m
}

func (c *Curd) Add(model interface{}, body interface{}) *Add {
	m := new(Add)
	m.initialize(c.db, model, nil)
	return m
}
