package service

import (
	"context"
	. "gopkg.in/check.v1"
)

func (x *MySuite) TestAdminFetch(c *C) {
	result, err := x.admin.Fetch(context.Background(), "1")
	if err != nil {
		c.Error(err)
	}
	c.Log(result)
}

func (x *MySuite) TestAdminClear(c *C) {
	err := x.acl.Clear(context.Background())
	if err != nil {
		c.Error(err)
	}
}
