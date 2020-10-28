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

func (c *Curd) Originlists(model interface{}, body OriginListsBody) *originListsModel {
	m := new(originListsModel)
	m.common = c.common
	m.model = model
	m.body = body
	return m
}

func (c *Curd) Lists(model interface{}, body ListsBody) *listsModel {
	m := new(listsModel)
	m.common = c.common
	m.model = model
	m.body = body
	return m
}

func (c *Curd) Get(model interface{}, body GetBody) *getModel {
	m := new(getModel)
	m.common = c.common
	m.model = model
	m.body = body
	return m
}

func (c *Curd) Add() *addModel {
	m := new(addModel)
	m.common = c.common
	return m
}

func (c *Curd) Edit(model interface{}, body EditBody) *editModel {
	m := new(editModel)
	m.common = c.common
	m.model = model
	m.body = body
	return m
}

func (c *Curd) Delete(model interface{}, body DeleteBody) *deleteModel {
	m := new(deleteModel)
	m.common = c.common
	m.model = model
	m.body = body
	return m
}

type Conditions [][]interface{}
type Query func(tx *gorm.DB) *gorm.DB
