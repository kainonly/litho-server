package service

import (
	"context"
	. "gopkg.in/check.v1"
)

func (x *MySuite) TestAclGet(c *C) {
	result, err := x.acl.Get(context.Background(), "resource", "1")
	if err != nil {
		c.Error(err)
	}
	c.Log(result)
}

func (x *MySuite) TestAclClear(c *C) {
	err := x.acl.Clear(context.Background())
	if err != nil {
		c.Error(err)
	}
}
