package service

import (
	. "gopkg.in/check.v1"
	"lab-api/bootstrap"
	"os"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type MySuite struct {
	lock     *Lock
	acl      *Acl
	resource *Resource
	role     *Role
	admin    *Admin
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
	dep := Dependent{
		Config: cfg,
		Db:     db,
		Redis:  redis,
	}
	x.lock = NewLock(&dep)
	x.acl = NewAcl(&dep)
	x.resource = NewResource(&dep)
	x.role = NewRole(&dep)
	x.admin = NewAdmin(&dep)
}
