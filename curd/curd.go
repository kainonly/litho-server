package curd

import (
	"gorm.io/gorm"
)

type Curd struct {
	common
}

type common struct {
	db *gorm.DB
}

func Initialize(db *gorm.DB) *Curd {
	c := new(Curd)
	c.db = db
	return c
}

func (c *Curd) Originlists(model interface{}, body OriginListsBody) *originLists {
	m := new(originLists)
	m.common = c.common
	m.model = model
	m.body = body
	return m
}

func (c *Curd) Lists(model interface{}, body ListsBody) *lists {
	m := new(lists)
	m.common = c.common
	m.model = model
	m.body = body
	return m
}

func (c *Curd) Get(model interface{}, body GetBody) *get {
	m := new(get)
	m.common = c.common
	m.model = model
	m.body = body
	return m
}

func (c *Curd) Add(model interface{}) *add {
	m := new(add)
	m.common = c.common
	m.model = model
	return m
}

func (c *Curd) Edit(model interface{}, body EditBody) *edit {
	m := new(edit)
	m.common = c.common
	m.model = model
	m.body = body
	return m
}

type Conditions [][]interface{}
type Query func(tx *gorm.DB)
