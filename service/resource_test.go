package service

import (
	"context"
	. "gopkg.in/check.v1"
)

func (x *MySuite) TestResourceFetch(c *C) {
	result, err := x.resource.Fetch(context.Background())
	if err != nil {
		c.Error(err)
	}
	c.Log(result)
}

func (x *MySuite) TestResourceClear(c *C) {
	err := x.acl.Clear(context.Background())
	if err != nil {
		c.Error(err)
	}
}
