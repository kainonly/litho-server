package curd

import (
	"log"
	"reflect"
)

func (c *Curd) Originlists(model interface{}) {
	slice := reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(model)), 0, 0)
	lists := reflect.New(slice.Type())
	c.db.Find(lists.Interface())
	log.Println(lists)
}
