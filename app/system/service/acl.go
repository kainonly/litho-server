package service

import (
	"context"
)

type Acl struct {
	*Dependency
	Key string
}

func NewAcl(d Dependency) *Acl {
	return &Acl{
		Dependency: &d,
		Key:        d.Config.RedisKey("acl"),
	}
}

func (x *Acl) Get(ctx context.Context) {

}

func (x *Acl) RefreshCache(ctx context.Context) (err error) {
	return
}

func (x *Acl) RemoveCache(ctx context.Context) error {
	return x.Redis.Del(ctx, x.Key).Err()
}
