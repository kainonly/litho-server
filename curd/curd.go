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

type common struct {
	db    *gorm.DB
	model interface{}
}

func (c *Curd) Originlists(model interface{}) *OriginLists {
	m := new(OriginLists)
	m.db = c.db
	m.model = model
	return m
}
