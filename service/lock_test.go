package service

import (
	"context"
	. "gopkg.in/check.v1"
)

func (x *MySuite) TestLocInc(c *C) {
	err := x.lock.Inc(context.Background(), "uid:1")
	if err != nil {
		c.Error(err)
	}
}
