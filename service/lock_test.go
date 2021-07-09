package service

import (
	"context"
	. "gopkg.in/check.v1"
)

func (x *MySuite) TestLockInc(c *C) {
	err := x.lock.Inc(context.Background(), "uid:1")
	if err != nil {
		c.Error(err)
	}
}

func (x *MySuite) TestLockRenew(c *C) {
	err := x.lock.Renew(context.Background(), "uid:1")
	if err != nil {
		c.Error(err)
	}
}

func (x *MySuite) TestLockCheck(c *C) {
	result, err := x.lock.Check(context.Background(), "uid:1")
	if err != nil {
		c.Error(err)
	}
	c.Log(result)
}

func (x *MySuite) TestLockCancel(c *C) {
	err := x.lock.Cancel(context.Background(), "uid:1")
	if err != nil {
		c.Error(err)
	}
}
