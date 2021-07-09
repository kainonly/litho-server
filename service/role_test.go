package service

import (
	"context"
	. "gopkg.in/check.v1"
)

func (x *MySuite) TestRoleFetch(c *C) {
	result, err := x.role.Fetch(context.Background(), []string{"*"}, "acl")
	if err != nil {
		c.Error(err)
	}
	c.Log(result)
	result, err = x.role.Fetch(context.Background(), []string{"*"}, "resource")
	if err != nil {
		c.Error(err)
	}
	c.Log(result)
}

func (x *MySuite) TestRoleClear(c *C) {
	err := x.acl.Clear(context.Background())
	if err != nil {
		c.Error(err)
	}
}
