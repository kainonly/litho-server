package curd

import (
	"gorm.io/gorm"
	"reflect"
	"van-api/helper/res"
)

type Curd struct {
	db *gorm.DB
}

func Initialize(db *gorm.DB) *Curd {
	c := new(Curd)
	c.db = db
	return c
}

func (c *Curd) Originlists(model interface{}) interface{} {
	slice := reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(model)), 0, 0)
	lists := reflect.New(slice.Type())
	c.db.Find(lists.Interface())
	return res.Data(lists.Interface())
}
