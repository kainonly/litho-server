package service

import (
	. "gopkg.in/check.v1"
	"lab-api/bootstrap"
	"os"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type MySuite struct {
	lock *Lock
	acl  *Acl
}

var _ = Suite(&MySuite{})

func (x *MySuite) SetUpTest(c *C) {
	os.Chdir("../")
	cfg, err := bootstrap.LoadConfiguration()
	if err != nil {
		c.Error(err)
	}
	db, err := bootstrap.InitializeDatabase(cfg)
	if err != nil {
		c.Error(err)
	}
	redis := bootstrap.InitializeRedis(cfg)
	x.lock = NewLock(cfg, redis)
	x.acl = NewAcl(cfg, db, redis)
}
